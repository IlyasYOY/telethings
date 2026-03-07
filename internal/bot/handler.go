package bot

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/IlyasYOY/telethings/internal/thingser"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const thingsListToday = "Today"
const thingsListInbox = "Inbox"
const thingsListAnytime = "Anytime"
const thingsListSomeday = "Someday"
const tasksPageSize = 10

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock@v3.4.7  -i MessageSender -o message_sender_mock_test.go -p bot_test -g

// MessageSender sends text replies to a Telegram chat.
type MessageSender interface {
	Send(chatID int64, text string) error
	SendWithInlineKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) error
	SendTyping(chatID int64) error
	AckCallback(callbackID string) error
}

// Handler dispatches Telegram updates to the appropriate command handler.
type Handler struct {
	sender         MessageSender
	opener         opener
	reader         thingsReader
	store          taskStore
	allowedUserIDs map[int64]bool
}

// NewHandler creates a Handler.
func NewHandler(sender MessageSender, o opener, r thingsReader, store taskStore, allowedUserIDs []int64) *Handler {
	idMap := make(map[int64]bool, len(allowedUserIDs))
	for _, id := range allowedUserIDs {
		idMap[id] = true
	}
	return &Handler{
		sender:         sender,
		opener:         o,
		reader:         r,
		store:          store,
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
	if msg.From == nil || !h.allowedUserIDs[msg.From.ID] {
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
	case "tags":
		return h.handleTags(msg)
	case "task":
		return h.handleTask(msg)
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

	if tag, ok := parseTagSelectionCallback(callback.Data); ok {
		if err := h.sender.AckCallback(callback.ID); err != nil {
			return err
		}
		return h.handlePaginatedTagTasks(callback.Message.Chat.ID, tag, 0)
	}
	if op, number, ok := parseTaskOperationCallback(callback.Data); ok {
		if err := h.sender.AckCallback(callback.ID); err != nil {
			return err
		}
		return h.handleTaskOperation(callback.Message.Chat.ID, number, op)
	}
	if tag, page, ok := parseTagPaginationCallback(callback.Data); ok {
		if err := h.sender.AckCallback(callback.ID); err != nil {
			return err
		}
		return h.handlePaginatedTagTasks(callback.Message.Chat.ID, tag, page)
	}

	list, page, ok := parsePaginationCallback(callback.Data)
	if ok {
		if err := h.sender.AckCallback(callback.ID); err != nil {
			return err
		}
		return h.handlePaginatedTaskList(callback.Message.Chat.ID, list, page, "📭 List is empty!")
	}

	return h.sender.AckCallback(callback.ID)
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

func parseTagSelectionCallback(data string) (tag string, ok bool) {
	const prefix = "tagsel:"
	if !strings.HasPrefix(data, prefix) {
		return "", false
	}
	tag, err := url.QueryUnescape(strings.TrimPrefix(data, prefix))
	if err != nil || strings.TrimSpace(tag) == "" {
		return "", false
	}
	return tag, true
}

func parseTagPaginationCallback(data string) (tag string, page int, ok bool) {
	const prefix = "tagpage:"
	if !strings.HasPrefix(data, prefix) {
		return "", 0, false
	}
	parts := strings.SplitN(strings.TrimPrefix(data, prefix), ":", 2)
	if len(parts) != 2 {
		return "", 0, false
	}
	tag, err := url.QueryUnescape(parts[0])
	if err != nil || strings.TrimSpace(tag) == "" {
		return "", 0, false
	}
	page, err = strconv.Atoi(parts[1])
	if err != nil || page < 0 {
		return "", 0, false
	}
	return tag, page, true
}

func tagSelectionCallbackData(tag string) string {
	return "tagsel:" + url.QueryEscape(tag)
}

func tagPageCallbackData(tag string, page int) string {
	return "tagpage:" + url.QueryEscape(tag) + ":" + strconv.Itoa(page)
}

func parseTaskOperationCallback(data string) (operation string, number int, ok bool) {
	const prefix = "taskop:"
	if !strings.HasPrefix(data, prefix) {
		return "", 0, false
	}
	parts := strings.SplitN(strings.TrimPrefix(data, prefix), ":", 2)
	if len(parts) != 2 {
		return "", 0, false
	}
	switch parts[0] {
	case "done", "undo", "cancel":
		operation = parts[0]
	default:
		return "", 0, false
	}
	number, err := strconv.Atoi(parts[1])
	if err != nil || number <= 0 {
		return "", 0, false
	}
	return operation, number, true
}

func taskOperationCallbackData(operation string, number int) string {
	return "taskop:" + operation + ":" + strconv.Itoa(number)
}

func (h *Handler) handlePaginatedTaskList(chatID int64, list string, page int, emptyMsg string) error {
	if err := h.sender.SendTyping(chatID); err != nil {
		return err
	}
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
	if err := h.saveTaskList(chatID, "list:"+strings.ToLower(list), page*tasksPageSize+1, tasks); err != nil {
		return err
	}

	text, keyboard := formatPaginatedTasks(list, tasks, page, hasNext)
	if len(keyboard.InlineKeyboard) == 0 {
		return h.sender.Send(chatID, text)
	}
	return h.sender.SendWithInlineKeyboard(chatID, text, keyboard)
}

func (h *Handler) handlePaginatedTagTasks(chatID int64, tag string, page int) error {
	if err := h.sender.SendTyping(chatID); err != nil {
		return err
	}
	if page < 0 {
		page = 0
	}

	offset := page * tasksPageSize
	tasks, err := h.reader.TasksByTagPage(tag, offset, tasksPageSize+1)
	if err != nil {
		return fmt.Errorf("read tag tasks: %w", err)
	}
	if len(tasks) == 0 {
		if page > 0 {
			return h.handlePaginatedTagTasks(chatID, tag, page-1)
		}
		return h.sender.Send(chatID, fmt.Sprintf("📭 No tasks found for tag: %s", tag))
	}

	hasNext := len(tasks) > tasksPageSize
	if hasNext {
		tasks = tasks[:tasksPageSize]
	}
	if err := h.saveTaskList(chatID, "tag:"+tag, page*tasksPageSize+1, tasks); err != nil {
		return err
	}

	text, keyboard := formatPaginatedTagTasks(tag, tasks, page, hasNext)
	if len(keyboard.InlineKeyboard) == 0 {
		return h.sender.Send(chatID, text)
	}
	return h.sender.SendWithInlineKeyboard(chatID, text, keyboard)
}

func formatPaginatedTasks(list string, tasks []thingser.Task, page int, hasNext bool) (string, tgbotapi.InlineKeyboardMarkup) {
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

func formatPaginatedTagTasks(tag string, tasks []thingser.Task, page int, hasNext bool) (string, tgbotapi.InlineKeyboardMarkup) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "🏷️ %s — page %d\n\n", tag, page+1)
	startNumber := page*tasksPageSize + 1
	for i, t := range tasks {
		fmt.Fprintf(&sb, "%d. %s\n", startNumber+i, formatTaskLine(t, true))
	}
	text := strings.TrimRight(sb.String(), "\n")

	var row []tgbotapi.InlineKeyboardButton
	if page > 0 {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev", tagPageCallbackData(tag, page-1)))
	}
	if hasNext {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("Next ➡️", tagPageCallbackData(tag, page+1)))
	}
	if len(row) == 0 {
		return text, tgbotapi.NewInlineKeyboardMarkup()
	}
	return text, tgbotapi.NewInlineKeyboardMarkup(row)
}

func (h *Handler) handleTaskList(msg *tgbotapi.Message, list, emptyMsg string) error {
	if err := h.sender.SendTyping(msg.Chat.ID); err != nil {
		return err
	}
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
		if err := h.saveTaskList(msg.Chat.ID, "list:"+strings.ToLower(list), 1, orderedTodayTasks(tasks)); err != nil {
			return err
		}
	} else {
		tasks = sortedTasksForDisplay(tasks)
		text = formatInboxTasks(tasks)
		if err := h.saveTaskList(msg.Chat.ID, "list:"+strings.ToLower(list), 1, tasks); err != nil {
			return err
		}
	}
	return h.sender.Send(msg.Chat.ID, text)
}

func (h *Handler) saveTaskList(chatID int64, scope string, startNumber int, tasks []thingser.Task) error {
	if h.store == nil {
		return nil
	}
	if err := h.store.SaveTaskList(chatID, scope, startNumber, tasks); err != nil {
		return fmt.Errorf("save task list mapping: %w", err)
	}
	return nil
}

func (h *Handler) handleTags(msg *tgbotapi.Message) error {
	if err := h.sender.SendTyping(msg.Chat.ID); err != nil {
		return err
	}
	tags, err := h.reader.Tags()
	if err != nil {
		return fmt.Errorf("read tags: %w", err)
	}
	if len(tags) == 0 {
		return h.sender.Send(msg.Chat.ID, "🏷️ No tags found!")
	}

	sort.Slice(tags, func(i, j int) bool {
		left := strings.ToLower(tags[i].Path)
		right := strings.ToLower(tags[j].Path)
		if left == right {
			return strings.ToLower(tags[i].Name) < strings.ToLower(tags[j].Name)
		}
		return left < right
	})
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, (len(tags)+1)/2)
	for i := 0; i < len(tags); i += 2 {
		row := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(tags[i].Path, tagSelectionCallbackData(tags[i].Path)),
		}
		if i+1 < len(tags) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(tags[i+1].Path, tagSelectionCallbackData(tags[i+1].Path)))
		}
		rows = append(rows, row)
	}

	return h.sender.SendWithInlineKeyboard(msg.Chat.ID, "🏷️ Choose a tag:", tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (h *Handler) handleTask(msg *tgbotapi.Message) error {
	if err := h.sender.SendTyping(msg.Chat.ID); err != nil {
		return err
	}
	number, err := strconv.Atoi(strings.TrimSpace(msg.CommandArguments()))
	if err != nil || number <= 0 {
		return h.sender.Send(msg.Chat.ID, "Usage: /task <number>")
	}
	if h.store == nil {
		return h.sender.Send(msg.Chat.ID, "Task storage is not configured.")
	}
	task, err := h.store.TaskByNumber(msg.Chat.ID, number)
	if err != nil {
		return h.sender.Send(msg.Chat.ID, "Task not found in the last shown list. Show a list first, then use /task <number>.")
	}
	text := formatTaskDetails(number, task)
	keyboard := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("✅ Done", taskOperationCallbackData("done", number)),
		tgbotapi.NewInlineKeyboardButtonData("↩️ Undo", taskOperationCallbackData("undo", number)),
		tgbotapi.NewInlineKeyboardButtonData("🚫 Cancel", taskOperationCallbackData("cancel", number)),
	})
	return h.sender.SendWithInlineKeyboard(msg.Chat.ID, text, keyboard)
}

func (h *Handler) handleTaskOperation(chatID int64, number int, operation string) error {
	if err := h.sender.SendTyping(chatID); err != nil {
		return err
	}
	if h.store == nil {
		return h.sender.Send(chatID, "Task storage is not configured.")
	}
	task, err := h.store.TaskByNumber(chatID, number)
	if err != nil {
		return h.sender.Send(chatID, "Task not found in the last shown list. Show a list first, then use /task <number>.")
	}
	if task.ID == "" {
		return h.sender.Send(chatID, "Task cannot be modified because its Things ID is missing.")
	}

	var successText string
	var updateErr error
	switch operation {
	case "done":
		successText = "✅ Task marked as done"
		updateErr = h.reader.SetTaskCompleted(task.ID, true)
	case "undo":
		successText = "↩️ Task marked as not completed"
		updateErr = h.reader.SetTaskCompleted(task.ID, false)
	case "cancel":
		successText = "🚫 Task canceled"
		updateErr = h.reader.SetTaskCanceled(task.ID, true)
	default:
		return h.sender.Send(chatID, "Unsupported task operation.")
	}
	if updateErr != nil {
		return fmt.Errorf("update task status: %w", updateErr)
	}
	return h.sender.Send(chatID, successText)
}

func formatTaskDetails(number int, task thingser.Task) string {
	var sb strings.Builder
	status := "open"
	if task.Canceled {
		status = "canceled"
	} else if task.Completed {
		status = "completed"
	}
	fmt.Fprintf(&sb, "🧩 Task #%d\n", number)
	fmt.Fprintf(&sb, "Title: %s\n", task.Title)
	fmt.Fprintf(&sb, "Status: %s\n", status)
	if task.Area != "" {
		fmt.Fprintf(&sb, "Area: %s\n", task.Area)
	}
	if task.Project != "" {
		fmt.Fprintf(&sb, "Project: %s\n", task.Project)
	}
	if task.Deadline != "" {
		fmt.Fprintf(&sb, "Deadline: %s\n", task.Deadline)
	}
	if len(task.Tags) > 0 {
		fmt.Fprintf(&sb, "Tags: %s\n", strings.Join(task.Tags, ","))
	}
	return strings.TrimRight(sb.String(), "\n")
}

func formatInboxTasks(tasks []thingser.Task) string {
	var sb strings.Builder
	for i, t := range tasks {
		fmt.Fprintf(&sb, "%d. %s\n", i+1, formatTaskLine(t, true))
	}
	return strings.TrimRight(sb.String(), "\n")
}

func sortedTasksForDisplay(tasks []thingser.Task) []thingser.Task {
	items := append([]thingser.Task(nil), tasks...)
	sort.Slice(items, func(i, j int) bool {
		return lessTaskForDisplay(items[i], items[j])
	})
	return items
}

func formatTodayTasks(tasks []thingser.Task) string {
	areaGroups := make(map[string][]thingser.Task)
	projectGroups := make(map[string][]thingser.Task)
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
	counter := 1
	for ai, area := range areas {
		if ai > 0 {
			sb.WriteString("\n\n")
		}
		fmt.Fprintf(&sb, "Area: %s\n", area)
		items := areaGroups[area]
		sort.Slice(items, func(i, j int) bool {
			return lessTaskForDisplay(items[i], items[j])
		})
		for _, task := range items {
			fmt.Fprintf(&sb, "  %d. %s\n", counter, formatTaskLine(task, false))
			counter++
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
			for _, task := range items {
				fmt.Fprintf(&sb, "  %d. %s\n", counter, formatTaskLine(task, false))
				counter++
			}
		}
	}

	return strings.TrimRight(sb.String(), "\n")
}

func orderedTodayTasks(tasks []thingser.Task) []thingser.Task {
	areaGroups := make(map[string][]thingser.Task)
	projectGroups := make(map[string][]thingser.Task)
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

	ordered := make([]thingser.Task, 0, len(tasks))
	for _, area := range areas {
		items := areaGroups[area]
		sort.Slice(items, func(i, j int) bool {
			return lessTaskForDisplay(items[i], items[j])
		})
		ordered = append(ordered, items...)
	}
	for _, project := range projects {
		items := projectGroups[project]
		sort.Slice(items, func(i, j int) bool {
			return lessTaskForDisplay(items[i], items[j])
		})
		ordered = append(ordered, items...)
	}
	return ordered
}

func formatTaskLine(task thingser.Task, includeAreaProject bool) string {
	prefix := "⬜ "
	if task.Canceled {
		prefix = "🚫 "
	} else if task.Completed {
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

func lessTaskForDisplay(a, b thingser.Task) bool {
	aClosed := a.Completed || a.Canceled
	bClosed := b.Completed || b.Canceled
	if aClosed != bClosed {
		return !aClosed && bClosed
	}
	return strings.ToLower(a.Title) < strings.ToLower(b.Title)
}

func (h *Handler) handleAdd(msg *tgbotapi.Message) error {
	args := strings.TrimSpace(msg.CommandArguments())
	input := parseAddCommandInput(args)
	if input == nil {
		return h.sender.Send(msg.Chat.ID, "Usage: /add <title> [when:<value>] [deadline:<value>] [tags:<csv>] [notes:<text>]")
	}

	task, err := h.reader.AddTask(thingser.AddTaskInput{
		Title:    input.title,
		When:     input.when,
		Deadline: input.deadline,
		Tags:     append([]string(nil), input.tags...),
		Notes:    input.notes,
	})
	if err != nil {
		return fmt.Errorf("add task: %w", err)
	}
	if h.store != nil {
		if err := h.store.SaveTaskList(msg.Chat.ID, "add:last", 1, []thingser.Task{task}); err != nil {
			return fmt.Errorf("save task mapping: %w", err)
		}
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
		"/someday - Show Someday tasks with pagination\n" +
		"/tags - Show all tags and read tasks by selected tag\n" +
		"/task <number> - Show task details and operation buttons\n"
	return h.sender.Send(msg.Chat.ID, text)
}
