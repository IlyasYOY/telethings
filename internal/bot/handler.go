package bot

import (
	"fmt"
	"sort"
	"strings"

	"github.com/IlyasYOY/telethings/internal/reader"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const thingsListToday = "Today"
const thingsListInbox = "Inbox"

// MessageSender sends text replies to a Telegram chat.
type MessageSender interface {
	Send(chatID int64, text string) error
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
	case "help":
		return h.handleHelp(msg)
	case "add":
		return h.handleAdd(msg)
	case "today":
		return h.handleTaskList(msg, thingsListToday, "📭 No tasks for today!")
	case "inbox":
		return h.handleTaskList(msg, thingsListInbox, "📭 Inbox is empty!")
	default:
		return h.sender.Send(msg.Chat.ID, "Unknown command. Use /help to see available commands.")
	}
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

func formatInboxTasks(tasks []reader.Task) string {
	var sb strings.Builder
	for i, t := range tasks {
		fmt.Fprintf(&sb, "%d. %s\n", i+1, formatTaskLine(t, true))
	}
	return strings.TrimRight(sb.String(), "\n")
}

func formatTodayTasks(tasks []reader.Task) string {
	groups := make(map[string]map[string][]reader.Task)
	for _, task := range tasks {
		area := task.Area
		if area == "" {
			area = "Other"
		}
		project := task.Project
		if project == "" {
			project = "Other"
		}
		if _, ok := groups[area]; !ok {
			groups[area] = make(map[string][]reader.Task)
		}
		groups[area][project] = append(groups[area][project], task)
	}

	areas := make([]string, 0, len(groups))
	for area := range groups {
		areas = append(areas, area)
	}
	sort.Strings(areas)

	var sb strings.Builder
	for ai, area := range areas {
		if ai > 0 {
			sb.WriteString("\n\n")
		}
		fmt.Fprintf(&sb, "Area: %s\n", area)

		projects := make([]string, 0, len(groups[area]))
		for project := range groups[area] {
			projects = append(projects, project)
		}
		sort.Strings(projects)

		for _, project := range projects {
			fmt.Fprintf(&sb, "  Project: %s\n", project)
			items := groups[area][project]
			sort.Slice(items, func(i, j int) bool {
				return strings.ToLower(items[i].Title) < strings.ToLower(items[j].Title)
			})
			for i, task := range items {
				fmt.Fprintf(&sb, "    %d. %s\n", i+1, formatTaskLine(task, false))
			}
		}
	}

	return strings.TrimRight(sb.String(), "\n")
}

func formatTaskLine(task reader.Task, includeAreaProject bool) string {
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
		return task.Title
	}
	return task.Title + " — " + strings.Join(parts, " | ")
}

func (h *Handler) handleAdd(msg *tgbotapi.Message) error {
	args := strings.TrimSpace(msg.CommandArguments())
	thingsURL := parseAddCommand(h.authToken, args)
	if thingsURL == "" {
		return h.sender.Send(msg.Chat.ID, "Usage: /add <title> [when:<value>] [tags:<csv>] [notes:<text>]")
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
		"  Options: [when:<value>] [tags:<csv>] [notes:<text>]\n\n" +
		"/today - Show today's tasks from Things 3\n" +
		"/inbox - Show your Things 3 inbox\n\n" +
		"/help - Show detailed command information\n"
	return h.sender.Send(msg.Chat.ID, text)
}

func (h *Handler) handleHelp(msg *tgbotapi.Message) error {
	text := "📚 Available Commands:\n\n" +
		"**/start** - Welcome message and quick help\n\n" +
		"**/add <title>** - Add a task to Things 3\n" +
		"  when:<value> - Schedule timing (e.g. today, next friday)\n" +
		"  tags:<csv> - Add tags (comma-separated)\n" +
		"  notes:<text> - Add detailed notes\n\n" +
		"**/today** - Show today's tasks from Things 3\n\n" +
		"**/inbox** - Show your Things 3 inbox\n\n" +
		"Examples:\n" +
		"  /add Buy milk\n" +
		"  /add Gym when:tomorrow tags:fitness\n" +
		"  /add Review notes:check email\n"
	return h.sender.Send(msg.Chat.ID, text)
}
