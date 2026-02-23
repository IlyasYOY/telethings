package bot

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/IlyasYOY/telethings/internal/thingsreader"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const thingsListToday = "Today"
const thingsListInbox = "Inbox"
const thingsListAnytime = "Anytime"
const thingsListSomeday = "Someday"
const tasksPageSize = 10

// MessageSender sends text replies to a Telegram chat.
type MessageSender interface {
	Send(chatID int64, text string) error
	SendWithInlineKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) error
	AckCallback(callbackID string) error
}

// Handler dispatches Telegram updates to the appropriate command handler.
type Handler struct {
	sender         MessageSender
	opener         opener
	reader         thingsReader
	authToken      string
	allowedUserIDs map[int64]bool
}

// NewHandler creates a Handler.
func NewHandler(sender MessageSender, o opener, r thingsReader, authToken string, allowedUserIDs []int64) *Handler {
	idMap := make(map[int64]bool, len(allowedUserIDs))
	for _, id := range allowedUserIDs {
		idMap[id] = true
	}
	return &Handler{
		sender:         sender,
		opener:         o,
		reader:         r,
		authToken:      authToken,
		allowedUserIDs: idMap,
	}
}

// Handle processes a single update.
func (h *Handler) Handle(update tgbotapi.Update) error {
	if update.CallbackQuery != nil {
		return h.handleCallback(update.CallbackQuery)
	}

	msg := update.Message
	if msg == nil || !msg.IsCommand() {
		return nil
	}

	// Check if user is allowed
	if !h.allowedUserIDs[msg.From.ID] {
		return nil
	}

	switch msg.Command() {
	case "start":
		return h.handleStart(msg)
	case "add":
		return h.handleAdd(msg)
	case "today":
		return h.handleTaskList(msg, thingsListToday, "📭 No tasks for today!")
	case "inbox":
		return h.handleTaskList(msg, thingsListInbox, "📭 Inbox is empty!")
	case "anytime":
		return h.handlePaginatedTaskList(msg.Chat.ID, thingsListAnytime, 0, "📭 Anytime is empty!")
	case "someday":
		return h.handlePaginatedTaskList(msg.Chat.ID, thingsListSomeday, 0, "📭 Someday is empty!")
	default:
		return h.sender.Send(msg.Chat.ID, "Unknown command. Use /start to see available commands.")
	}
}

func (h *Handler) handleCallback(callback *tgbotapi.CallbackQuery) error {
	if callback == nil || callback.Message == nil {
		return nil
	}
	if !h.allowedUserIDs[callback.From.ID] {
		return nil
	}

	list, page, ok := parsePaginationCallback(callback.Data)
	if !ok {
		return h.sender.AckCallback(callback.ID)
	}

	if err := h.sender.AckCallback(callback.ID); err != nil {
		return err
	}

	return h.handlePaginatedTaskList(callback.Message.Chat.ID, list, page, "📭 List is empty!")
}

func parsePaginationCallback(data string) (list string, page int, ok bool) {
	parts := strings.Split(data, ":")
	if len(parts) != 3 || parts[0] != "page" {
		return "", 0, false
	}

	page, err := strconv.Atoi(parts[2])
	if err != nil || page < 0 {
		return "", 0, false
	}

	switch parts[1] {
	case "anytime":
		return thingsListAnytime, page, true
	case "someday":
		return thingsListSomeday, page, true
	default:
		return "", 0, false
	}
}

func callbackListID(list string) string {
	switch list {
	case thingsListAnytime:
		return "anytime"
	case thingsListSomeday:
		return "someday"
	default:
		return ""
	}
}

func (h *Handler) handlePaginatedTaskList(chatID int64, list string, page int, emptyMsg string) error {
	if page < 0 {
		page = 0
	}

	offset := page * tasksPageSize
	tasks, err := h.reader.TasksInListPage(list, offset, tasksPageSize+1)
	if err != nil {
		return fmt.Errorf("read tasks: %w", err)
	}
	if len(tasks) == 0 {
		if page > 0 {
			return h.handlePaginatedTaskList(chatID, list, page-1, emptyMsg)
		}
		return h.sender.Send(chatID, emptyMsg)
	}

	hasNext := len(tasks) > tasksPageSize
	if hasNext {
		tasks = tasks[:tasksPageSize]
	}

	text, keyboard := formatPaginatedTasks(list, tasks, page, hasNext)
	return h.sender.SendWithInlineKeyboard(chatID, text, keyboard)
}

func formatPaginatedTasks(list string, tasks []thingsreader.Task, page int, hasNext bool) (string, tgbotapi.InlineKeyboardMarkup) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "📋 %s — page %d\n\n", list, page+1)
	startNumber := page*tasksPageSize + 1
	for i, t := range tasks {
		fmt.Fprintf(&sb, "%d. %s\n", startNumber+i, formatTaskLine(t, true))
	}
	text := strings.TrimRight(sb.String(), "\n")

	listID := callbackListID(list)
	var row []tgbotapi.InlineKeyboardButton
	if page > 0 {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev", fmt.Sprintf("page:%s:%d", listID, page-1)))
	}
	if hasNext {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("Next ➡️", fmt.Sprintf("page:%s:%d", listID, page+1)))
	}
	if len(row) == 0 {
		return text, tgbotapi.NewInlineKeyboardMarkup()
	}
	return text, tgbotapi.NewInlineKeyboardMarkup(row)
}

func (h *Handler) handleTaskList(msg *tgbotapi.Message, list, emptyMsg string) error {
	tasks, err := h.reader.TasksInList(list)
	if err != nil {
		return fmt.Errorf("read tasks: %w", err)
	}
	if len(tasks) == 0 {
		return h.sender.Send(msg.Chat.ID, emptyMsg)
	}

	var text string
	if list == thingsListToday {
		text = formatTodayTasks(tasks)
	} else {
		text = formatInboxTasks(tasks)
	}
	return h.sender.Send(msg.Chat.ID, text)
}

func formatInboxTasks(tasks []thingsreader.Task) string {
	items := append([]thingsreader.Task(nil), tasks...)
	sort.Slice(items, func(i, j int) bool {
		return lessTaskForDisplay(items[i], items[j])
	})

	var sb strings.Builder
	for i, t := range items {
		fmt.Fprintf(&sb, "%d. %s\n", i+1, formatTaskLine(t, true))
	}
	return strings.TrimRight(sb.String(), "\n")
}

func formatTodayTasks(tasks []thingsreader.Task) string {
	areaGroups := make(map[string][]thingsreader.Task)
	projectGroups := make(map[string][]thingsreader.Task)
	for _, task := range tasks {
		switch {
		case task.Area != "":
			areaGroups[task.Area] = append(areaGroups[task.Area], task)
		case task.Project != "":
			projectGroups[task.Project] = append(projectGroups[task.Project], task)
		default:
			areaGroups["Other"] = append(areaGroups["Other"], task)
		}
	}

	areas := make([]string, 0, len(areaGroups))
	for area := range areaGroups {
		areas = append(areas, area)
	}
	sort.Strings(areas)

	projects := make([]string, 0, len(projectGroups))
	for project := range projectGroups {
		projects = append(projects, project)
	}
	sort.Strings(projects)

	var sb strings.Builder
	for ai, area := range areas {
		if ai > 0 {
			sb.WriteString("\n\n")
		}
		fmt.Fprintf(&sb, "Area: %s\n", area)
		items := areaGroups[area]
		sort.Slice(items, func(i, j int) bool {
			return lessTaskForDisplay(items[i], items[j])
		})
		for i, task := range items {
			fmt.Fprintf(&sb, "  %d. %s\n", i+1, formatTaskLine(task, false))
		}
	}

	if len(projects) > 0 {
		if sb.Len() > 0 {
			sb.WriteString("\n\n")
		}
		for pi, project := range projects {
			if pi > 0 {
				sb.WriteString("\n\n")
			}
			fmt.Fprintf(&sb, "Project: %s\n", project)
			items := projectGroups[project]
			sort.Slice(items, func(i, j int) bool {
				return lessTaskForDisplay(items[i], items[j])
			})
			for i, task := range items {
				fmt.Fprintf(&sb, "  %d. %s\n", i+1, formatTaskLine(task, false))
			}
		}
	}

	return strings.TrimRight(sb.String(), "\n")
}

func formatTaskLine(task thingsreader.Task, includeAreaProject bool) string {
	prefix := "⬜ "
	if task.Completed {
		prefix = "✅ "
	}

	parts := make([]string, 0, 3)
	if includeAreaProject {
		switch {
		case task.Area != "" && task.Project != "":
			parts = append(parts, task.Area+"/"+task.Project)
		case task.Area != "":
			parts = append(parts, task.Area)
		case task.Project != "":
			parts = append(parts, task.Project)
		}
	}
	if task.Deadline != "" {
		parts = append(parts, "deadline:"+task.Deadline)
	}
	if len(task.Tags) > 0 {
		parts = append(parts, "tags:"+strings.Join(task.Tags, ","))
	}
	if len(parts) == 0 {
		return prefix + task.Title
	}
	return prefix + task.Title + " — " + strings.Join(parts, " | ")
}

func lessTaskForDisplay(a, b thingsreader.Task) bool {
	if a.Completed != b.Completed {
		return !a.Completed && b.Completed
	}
	return strings.ToLower(a.Title) < strings.ToLower(b.Title)
}

func (h *Handler) handleAdd(msg *tgbotapi.Message) error {
	args := strings.TrimSpace(msg.CommandArguments())
	thingsURL := parseAddCommand(h.authToken, args)
	if thingsURL == "" {
		return h.sender.Send(msg.Chat.ID, "Usage: /add <title> [when:<value>] [deadline:<value>] [tags:<csv>] [notes:<text>]")
	}

	if err := h.opener.Open(thingsURL); err != nil {
		return fmt.Errorf("open things URL: %w", err)
	}

	return h.sender.Send(msg.Chat.ID, "✅ Added to Things3")
}

func (h *Handler) handleStart(msg *tgbotapi.Message) error {
	text := "👋 Welcome to Telethings!\n\n" +
		"A Telegram bot that integrates with Things 3 task management.\n\n" +
		"📋 Available commands:\n\n" +
		"/add <title> - Add a task to Things 3\n" +
		"  Options: [when:<value>] [deadline:<value>] [tags:<csv>] [notes:<text>]\n\n" +
		"/today - Show today's tasks from Things 3\n" +
		"/inbox - Show your Things 3 inbox\n" +
		"/anytime - Show Anytime tasks with pagination\n" +
		"/someday - Show Someday tasks with pagination\n"
	return h.sender.Send(msg.Chat.ID, text)
}
