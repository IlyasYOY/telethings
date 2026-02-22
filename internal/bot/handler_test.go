package bot_test

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/IlyasYOY/telethings/internal/bot"
	"github.com/IlyasYOY/telethings/internal/opener/openertest"
)

// fakeSender records sent messages without making network calls.
type fakeSender struct {
	messages []struct {
		chatID int64
		text   string
	}
}

func (f *fakeSender) Send(chatID int64, text string) error {
	f.messages = append(f.messages, struct {
		chatID int64
		text   string
	}{chatID, text})
	return nil
}

func newTestUpdate(chatID int64, text string) tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: chatID},
			Text: text,
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: len(commandFrom(text))},
			},
		},
	}
}

// commandFrom extracts the /command portion (up to the first space).
func commandFrom(text string) string {
	for i, ch := range text {
		if ch == ' ' {
			return text[:i]
		}
	}
	return text
}

func TestHandler_HandleAdd_ValidCommand(t *testing.T) {
	const authToken = "tok"
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender := &fakeSender{}
	h := bot.NewHandler(sender, rec, authToken)

	update := newTestUpdate(chatID, "/add Buy milk")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rec.URLs) != 1 {
		t.Fatalf("expected 1 URL opened, got %d", len(rec.URLs))
	}
	if len(sender.messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(sender.messages))
	}
	if sender.messages[0].chatID != chatID {
		t.Errorf("reply chatID = %d, want %d", sender.messages[0].chatID, chatID)
	}
}

func TestHandler_HandleAdd_EmptyCommand(t *testing.T) {
	const authToken = "tok"
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender := &fakeSender{}
	h := bot.NewHandler(sender, rec, authToken)

	update := newTestUpdate(chatID, "/add")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rec.URLs) != 0 {
		t.Errorf("expected 0 URLs opened, got %d", len(rec.URLs))
	}
	if len(sender.messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(sender.messages))
	}
}

func TestHandler_HandleAdd_UnknownCommand(t *testing.T) {
	const authToken = "tok"
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender := &fakeSender{}
	h := bot.NewHandler(sender, rec, authToken)

	update := newTestUpdate(chatID, "/unknown")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rec.URLs) != 0 {
		t.Errorf("expected 0 URLs opened, got %d", len(rec.URLs))
	}
	if len(sender.messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(sender.messages))
	}
}

func TestHandler_HandleAdd_NonCommandMessage(t *testing.T) {
	const authToken = "tok"
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender := &fakeSender{}
	h := bot.NewHandler(sender, rec, authToken)

	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: chatID},
			Text: "just text",
		},
	}
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rec.URLs) != 0 || len(sender.messages) != 0 {
		t.Error("expected no action for non-command message")
	}
}
