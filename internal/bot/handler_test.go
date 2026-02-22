package bot_test

import (
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/IlyasYOY/telethings/internal/bot"
	"github.com/IlyasYOY/telethings/internal/opener/openertest"
	"github.com/IlyasYOY/telethings/internal/reader/readertest"
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

func newTestUpdate(userID, chatID int64, text string) tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: userID},
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
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender := &fakeSender{}
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/add Buy milk")
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
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender := &fakeSender{}
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/add")
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
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender := &fakeSender{}
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/unknown")
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
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender := &fakeSender{}
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, authToken, []int64{userID})

	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: userID},
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

func TestHandler_UnauthorizedUser(t *testing.T) {
	const authToken = "tok"
	const allowedUserID = int64(42)
	const unauthorizedUserID = int64(999)
	const chatID = int64(999)

	rec := &openertest.RecordingOpener{}
	sender := &fakeSender{}
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, authToken, []int64{allowedUserID})

	update := newTestUpdate(unauthorizedUserID, chatID, "/add Buy milk")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rec.URLs) != 0 {
		t.Errorf("expected 0 URLs opened for unauthorized user, got %d", len(rec.URLs))
	}
	if len(sender.messages) != 0 {
		t.Errorf("expected 0 replies for unauthorized user, got %d", len(sender.messages))
	}
}

func TestHandler_MultipleAllowedUsers(t *testing.T) {
	const authToken = "tok"
	allowedUserIDs := []int64{42, 123, 456}
	const chatID = int64(123)

	rec := &openertest.RecordingOpener{}
	sender := &fakeSender{}
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, authToken, allowedUserIDs)

	for _, userID := range allowedUserIDs {
		rec.URLs = []string{}
		sender.messages = []struct {
			chatID int64
			text   string
		}{}

		update := newTestUpdate(userID, chatID, "/add Test")
		if err := h.Handle(update); err != nil {
			t.Fatalf("unexpected error for user %d: %v", userID, err)
		}

		if len(rec.URLs) != 1 {
			t.Errorf("expected 1 URL for user %d, got %d", userID, len(rec.URLs))
		}
	}
}

func TestHandler_HandleToday_WithTasks(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender := &fakeSender{}
	rdr := &readertest.RecordingReader{Tasks: []string{"Buy milk", "Call dentist"}}
	h := bot.NewHandler(sender, rec, rdr, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/today")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sender.messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(sender.messages))
	}
	reply := sender.messages[0].text
	if !strings.Contains(reply, "Buy milk") || !strings.Contains(reply, "Call dentist") {
		t.Errorf("reply missing expected tasks: %q", reply)
	}
}

func TestHandler_HandleToday_Empty(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender := &fakeSender{}
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/today")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sender.messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(sender.messages))
	}
	if !strings.Contains(sender.messages[0].text, "No tasks for today") {
		t.Errorf("expected empty-list message, got: %q", sender.messages[0].text)
	}
}

func TestHandler_HandleInbox_WithTasks(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender := &fakeSender{}
	rdr := &readertest.RecordingReader{Tasks: []string{"Read book", "Fix bug"}}
	h := bot.NewHandler(sender, rec, rdr, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/inbox")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sender.messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(sender.messages))
	}
	reply := sender.messages[0].text
	if !strings.Contains(reply, "Read book") || !strings.Contains(reply, "Fix bug") {
		t.Errorf("reply missing expected tasks: %q", reply)
	}
}

func TestHandler_HandleInbox_Empty(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender := &fakeSender{}
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/inbox")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sender.messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(sender.messages))
	}
	if !strings.Contains(sender.messages[0].text, "Inbox is empty") {
		t.Errorf("expected empty-list message, got: %q", sender.messages[0].text)
	}
}
