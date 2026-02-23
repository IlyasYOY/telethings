package bot_test

import (
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/IlyasYOY/telethings/internal/bot"
	"github.com/IlyasYOY/telethings/internal/opener/openertest"
	"github.com/IlyasYOY/telethings/internal/reader"
	"github.com/IlyasYOY/telethings/internal/reader/readertest"
)

type sentMessage struct {
	chatID       int64
	text         string
	withKeyboard bool
	keyboard     tgbotapi.InlineKeyboardMarkup
}

type fakeSender struct {
	messages      []sentMessage
	ackCallbacks  []string
	sendErr       error
	sendInlineErr error
	ackErr        error
}

func (s *fakeSender) Send(chatID int64, text string) error {
	if s.sendErr != nil {
		return s.sendErr
	}
	s.messages = append(s.messages, sentMessage{chatID: chatID, text: text})
	return nil
}

func (s *fakeSender) SendWithInlineKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	if s.sendInlineErr != nil {
		return s.sendInlineErr
	}
	s.messages = append(s.messages, sentMessage{chatID: chatID, text: text, withKeyboard: true, keyboard: keyboard})
	return nil
}

func (s *fakeSender) AckCallback(callbackID string) error {
	if s.ackErr != nil {
		return s.ackErr
	}
	s.ackCallbacks = append(s.ackCallbacks, callbackID)
	return nil
}

func newSenderMock() (*fakeSender, *[]sentMessage) {
	sender := &fakeSender{}
	return sender, &sender.messages
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

func newCallbackUpdate(userID, chatID int64, callbackID, data string) tgbotapi.Update {
	return tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			ID:   callbackID,
			Data: data,
			From: &tgbotapi.User{ID: userID},
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: chatID},
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
	sender, messages := newSenderMock()
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/add Buy milk")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rec.URLs) != 1 {
		t.Fatalf("expected 1 URL opened, got %d", len(rec.URLs))
	}
	if len(*messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(*messages))
	}
	if (*messages)[0].chatID != chatID {
		t.Errorf("reply chatID = %d, want %d", (*messages)[0].chatID, chatID)
	}
}

func TestHandler_HandleAdd_EmptyCommand(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/add")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rec.URLs) != 0 {
		t.Errorf("expected 0 URLs opened, got %d", len(rec.URLs))
	}
	if len(*messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(*messages))
	}
}

func TestHandler_HandleAdd_UnknownCommand(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/unknown")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rec.URLs) != 0 {
		t.Errorf("expected 0 URLs opened, got %d", len(rec.URLs))
	}
	if len(*messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(*messages))
	}
	if !strings.Contains((*messages)[0].text, "/start") {
		t.Fatalf("expected unknown-command hint to mention /start, got %q", (*messages)[0].text)
	}
}

func TestHandler_HandleAdd_NonCommandMessage(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
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

	if len(rec.URLs) != 0 || len(*messages) != 0 {
		t.Error("expected no action for non-command message")
	}
}

func TestHandler_UnauthorizedUser(t *testing.T) {
	const authToken = "tok"
	const allowedUserID = int64(42)
	const unauthorizedUserID = int64(999)
	const chatID = int64(999)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, authToken, []int64{allowedUserID})

	update := newTestUpdate(unauthorizedUserID, chatID, "/add Buy milk")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rec.URLs) != 0 {
		t.Errorf("expected 0 URLs opened for unauthorized user, got %d", len(rec.URLs))
	}
	if len(*messages) != 0 {
		t.Errorf("expected 0 replies for unauthorized user, got %d", len(*messages))
	}
}

func TestHandler_MultipleAllowedUsers(t *testing.T) {
	const authToken = "tok"
	allowedUserIDs := []int64{42, 123, 456}
	const chatID = int64(123)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, authToken, allowedUserIDs)

	for _, userID := range allowedUserIDs {
		rec.URLs = []string{}
		*messages = nil

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
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{Tasks: []reader.Task{{Title: "Buy milk"}, {Title: "Call dentist"}}}
	h := bot.NewHandler(sender, rec, rdr, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/today")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(*messages))
	}
	reply := (*messages)[0].text
	if !strings.Contains(reply, "⬜ Buy milk") || !strings.Contains(reply, "⬜ Call dentist") {
		t.Errorf("reply missing expected tasks: %q", reply)
	}
}

func TestHandler_HandleToday_Empty(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/today")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(*messages))
	}
	if !strings.Contains((*messages)[0].text, "No tasks for today") {
		t.Errorf("expected empty-list message, got: %q", (*messages)[0].text)
	}
}

func TestHandler_HandleInbox_WithTasks(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{Tasks: []reader.Task{{Title: "Read book"}, {Title: "Fix bug"}}}
	h := bot.NewHandler(sender, rec, rdr, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/inbox")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(*messages))
	}
	reply := (*messages)[0].text
	if !strings.Contains(reply, "⬜ Read book") || !strings.Contains(reply, "⬜ Fix bug") {
		t.Errorf("reply missing expected tasks: %q", reply)
	}
}

func TestHandler_HandleInbox_Empty(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/inbox")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(*messages))
	}
	if !strings.Contains((*messages)[0].text, "Inbox is empty") {
		t.Errorf("expected empty-list message, got: %q", (*messages)[0].text)
	}
}

func TestHandler_HandleInbox_WithMetadata(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{
		Tasks: []reader.Task{
			{Title: "Read book", Area: "Life", Project: "Reading", Deadline: "Friday", Tags: []string{"home", "fun"}},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/inbox")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(*messages))
	}
	reply := (*messages)[0].text
	if !strings.Contains(reply, "⬜ Read book — Life/Reading | deadline:Friday | tags:home,fun") {
		t.Errorf("unexpected inbox format: %q", reply)
	}
}

func TestHandler_HandleToday_GroupedByAreaThenProject(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{
		Tasks: []reader.Task{
			{Title: "Task A", Area: "Work"},
			{Title: "Task B", Project: "Alpha"},
			{Title: "Task C", Area: "Life"},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/today")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(*messages))
	}
	reply := (*messages)[0].text
	if !strings.Contains(reply, "Area: Life") || !strings.Contains(reply, "Area: Work") {
		t.Errorf("expected area sections, got: %q", reply)
	}
	if !strings.Contains(reply, "Project: Alpha") {
		t.Errorf("expected project section, got: %q", reply)
	}
	if strings.Index(reply, "Area: Work") > strings.Index(reply, "Project: Alpha") {
		t.Errorf("expected areas before projects, got: %q", reply)
	}
}

func TestHandler_HandleInbox_CompletedAtBottom(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{
		Tasks: []reader.Task{
			{Title: "Done task", Completed: true},
			{Title: "Open task"},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/inbox")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reply := (*messages)[0].text
	if strings.Index(reply, "⬜ Open task") > strings.Index(reply, "✅ Done task") {
		t.Errorf("expected open task before completed task, got: %q", reply)
	}
}

func TestHandler_HandleToday_CompletedAtBottomInsideSection(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{
		Tasks: []reader.Task{
			{Title: "Done task", Area: "Work", Completed: true},
			{Title: "Open task", Area: "Work"},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/today")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reply := (*messages)[0].text
	if strings.Index(reply, "⬜ Open task") > strings.Index(reply, "✅ Done task") {
		t.Errorf("expected open task before completed task in area section, got: %q", reply)
	}
}

func TestHandler_HandleAnytime_WithPagination(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{
		Tasks: []reader.Task{
			{Title: "Task 01"}, {Title: "Task 02"}, {Title: "Task 03"}, {Title: "Task 04"}, {Title: "Task 05"},
			{Title: "Task 06"}, {Title: "Task 07"}, {Title: "Task 08"}, {Title: "Task 09"}, {Title: "Task 10"},
			{Title: "Task 11"},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/anytime")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(*messages))
	}
	reply := (*messages)[0]
	if !reply.withKeyboard {
		t.Fatal("expected inline keyboard")
	}
	if !strings.Contains(reply.text, "Anytime — page 1") {
		t.Fatalf("expected first page header, got %q", reply.text)
	}
	if rdr.LastPageList != "Anytime" || rdr.LastPageOffset != 0 || rdr.LastPageLimit != 11 {
		t.Fatalf("expected paged reader call with offset=0 limit=11, got list=%q offset=%d limit=%d", rdr.LastPageList, rdr.LastPageOffset, rdr.LastPageLimit)
	}
	if len(reply.keyboard.InlineKeyboard) != 1 || len(reply.keyboard.InlineKeyboard[0]) != 1 {
		t.Fatalf("expected one Next button, got %#v", reply.keyboard.InlineKeyboard)
	}
	if reply.keyboard.InlineKeyboard[0][0].CallbackData == nil || *reply.keyboard.InlineKeyboard[0][0].CallbackData != "page:anytime:1" {
		t.Fatalf("unexpected callback data: %#v", reply.keyboard.InlineKeyboard[0][0].CallbackData)
	}
}

func TestHandler_CallbackPagination_NextPage(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{
		Tasks: []reader.Task{
			{Title: "Task 01"}, {Title: "Task 02"}, {Title: "Task 03"}, {Title: "Task 04"}, {Title: "Task 05"},
			{Title: "Task 06"}, {Title: "Task 07"}, {Title: "Task 08"}, {Title: "Task 09"}, {Title: "Task 10"},
			{Title: "Task 11"},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, authToken, []int64{userID})

	update := newCallbackUpdate(userID, chatID, "cb-1", "page:anytime:1")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sender.ackCallbacks) != 1 || sender.ackCallbacks[0] != "cb-1" {
		t.Fatalf("expected callback ack, got %#v", sender.ackCallbacks)
	}
	if len(*messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(*messages))
	}
	reply := (*messages)[0]
	if !strings.Contains(reply.text, "Anytime — page 2") || !strings.Contains(reply.text, "11. ⬜ Task 11") {
		t.Fatalf("expected page 2 with task 11, got %q", reply.text)
	}
	if rdr.LastPageList != "Anytime" || rdr.LastPageOffset != 10 || rdr.LastPageLimit != 11 {
		t.Fatalf("expected paged reader call with offset=10 limit=11, got list=%q offset=%d limit=%d", rdr.LastPageList, rdr.LastPageOffset, rdr.LastPageLimit)
	}
	if len(reply.keyboard.InlineKeyboard) != 1 || len(reply.keyboard.InlineKeyboard[0]) != 1 {
		t.Fatalf("expected one Prev button, got %#v", reply.keyboard.InlineKeyboard)
	}
	if reply.keyboard.InlineKeyboard[0][0].CallbackData == nil || *reply.keyboard.InlineKeyboard[0][0].CallbackData != "page:anytime:0" {
		t.Fatalf("unexpected callback data: %#v", reply.keyboard.InlineKeyboard[0][0].CallbackData)
	}
}

func TestHandler_HandleSomeday_Empty(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, authToken, []int64{userID})

	update := newTestUpdate(userID, chatID, "/someday")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(*messages))
	}
	if !strings.Contains((*messages)[0].text, "Someday is empty") {
		t.Fatalf("expected empty Someday message, got %q", (*messages)[0].text)
	}
}
