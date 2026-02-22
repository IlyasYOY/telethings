package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MessageSender sends text replies to a Telegram chat.
type MessageSender interface {
	Send(chatID int64, text string) error
}

// Handler dispatches Telegram updates to the appropriate command handler.
type Handler struct {
	sender         MessageSender
	opener         opener
	authToken      string
	allowedUserIDs map[int64]bool
}

// NewHandler creates a Handler.
func NewHandler(sender MessageSender, o opener, authToken string, allowedUserIDs []int64) *Handler {
	idMap := make(map[int64]bool, len(allowedUserIDs))
	for _, id := range allowedUserIDs {
		idMap[id] = true
	}
	return &Handler{
		sender:         sender,
		opener:         o,
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
	default:
		return h.sender.Send(msg.Chat.ID, "Unknown command. Use /help to see available commands.")
	}
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
		"Examples:\n" +
		"  /add Buy milk\n" +
		"  /add Gym when:tomorrow tags:fitness\n" +
		"  /add Review notes:check email\n"
	return h.sender.Send(msg.Chat.ID, text)
}
