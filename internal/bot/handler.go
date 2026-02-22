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
	sender    MessageSender
	opener    opener
	authToken string
}

// NewHandler creates a Handler.
func NewHandler(sender MessageSender, o opener, authToken string) *Handler {
	return &Handler{sender: sender, opener: o, authToken: authToken}
}

// Handle processes a single update.
func (h *Handler) Handle(update tgbotapi.Update) error {
	msg := update.Message
	if msg == nil || !msg.IsCommand() {
		return nil
	}

	switch msg.Command() {
	case "add":
		return h.handleAdd(msg)
	default:
		return h.sender.Send(msg.Chat.ID, "Unknown command. Try /add <title>")
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
