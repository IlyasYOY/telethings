package bot_test

import (
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/IlyasYOY/telethings/internal/bot"
	"github.com/IlyasYOY/telethings/internal/opener/openertest"
	"github.com/IlyasYOY/telethings/internal/thingsreader"
	"github.com/IlyasYOY/telethings/internal/thingsreader/readertest"
)

type sentMessage struct {
	chatID       int64
	text         string
	withKeyboard bool
	keyboard     tgbotapi.InlineKeyboardMarkup
}

func newSenderMock() (*RecordingSender, *[]sentMessage) {
	sender, messages := NewRecordingSender()
	return sender.(*RecordingSender), messages
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
	rdr := &readertest.RecordingReader{}
	h := bot.NewHandler(sender, rec, rdr, nil, []int64{userID})

	update := newTestUpdate(userID, chatID, "/add Buy milk")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rdr.LastAddInput.Title != "Buy milk" {
		t.Fatalf("expected AddTask to be called with title %q, got %#v", "Buy milk", rdr.LastAddInput)
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
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, nil, []int64{userID})

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
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, nil, []int64{userID})

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
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, nil, []int64{userID})

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
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, nil, []int64{allowedUserID})

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
	rdr := &readertest.RecordingReader{}
	h := bot.NewHandler(sender, rec, rdr, nil, allowedUserIDs)

	for _, userID := range allowedUserIDs {
		rec.URLs = []string{}
		*messages = nil

		update := newTestUpdate(userID, chatID, "/add Test")
		if err := h.Handle(update); err != nil {
			t.Fatalf("unexpected error for user %d: %v", userID, err)
		}

		if rdr.LastAddInput.Title != "Test" {
			t.Errorf("expected add input title %q for user %d, got %#v", "Test", userID, rdr.LastAddInput)
		}
	}
}

func TestHandler_HandleToday_WithTasks(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{Tasks: []thingsreader.Task{{Title: "Buy milk"}, {Title: "Call dentist"}}}
	h := bot.NewHandler(sender, rec, rdr, nil, []int64{userID})

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
	if len(sender.typingChats) != 1 || sender.typingChats[0] != chatID {
		t.Fatalf("expected typing for chat %d, got %#v", chatID, sender.typingChats)
	}
}

func TestHandler_HandleToday_Empty(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, nil, []int64{userID})

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
	rdr := &readertest.RecordingReader{Tasks: []thingsreader.Task{{Title: "Read book"}, {Title: "Fix bug"}}}
	h := bot.NewHandler(sender, rec, rdr, nil, []int64{userID})

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
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, nil, []int64{userID})

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
		Tasks: []thingsreader.Task{
			{Title: "Read book", Area: "Life", Project: "Reading", Deadline: "Friday", Tags: []string{"home", "fun"}},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, nil, []int64{userID})

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
		Tasks: []thingsreader.Task{
			{Title: "Task A", Area: "Work"},
			{Title: "Task B", Project: "Alpha"},
			{Title: "Task C", Area: "Life"},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, nil, []int64{userID})

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
		Tasks: []thingsreader.Task{
			{Title: "Done task", Completed: true},
			{Title: "Open task"},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, nil, []int64{userID})

	update := newTestUpdate(userID, chatID, "/inbox")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reply := (*messages)[0].text
	if strings.Index(reply, "⬜ Open task") > strings.Index(reply, "✅ Done task") {
		t.Errorf("expected open task before completed task, got: %q", reply)
	}
}

func TestHandler_HandleInbox_CanceledShownWithCancelSymbol(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{
		Tasks: []thingsreader.Task{
			{Title: "Canceled task", Canceled: true},
			{Title: "Open task"},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, nil, []int64{userID})

	update := newTestUpdate(userID, chatID, "/inbox")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reply := (*messages)[0].text
	if !strings.Contains(reply, "🚫 Canceled task") {
		t.Fatalf("expected canceled symbol for canceled task, got: %q", reply)
	}
}

func TestHandler_HandleToday_CompletedAtBottomInsideSection(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{
		Tasks: []thingsreader.Task{
			{Title: "Done task", Area: "Work", Completed: true},
			{Title: "Open task", Area: "Work"},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, nil, []int64{userID})

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
		Tasks: []thingsreader.Task{
			{Title: "Task 01"}, {Title: "Task 02"}, {Title: "Task 03"}, {Title: "Task 04"}, {Title: "Task 05"},
			{Title: "Task 06"}, {Title: "Task 07"}, {Title: "Task 08"}, {Title: "Task 09"}, {Title: "Task 10"},
			{Title: "Task 11"},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, nil, []int64{userID})

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
		Tasks: []thingsreader.Task{
			{Title: "Task 01"}, {Title: "Task 02"}, {Title: "Task 03"}, {Title: "Task 04"}, {Title: "Task 05"},
			{Title: "Task 06"}, {Title: "Task 07"}, {Title: "Task 08"}, {Title: "Task 09"}, {Title: "Task 10"},
			{Title: "Task 11"},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, nil, []int64{userID})

	update := newCallbackUpdate(userID, chatID, "cb-1", "page:anytime:1")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sender.ackCallbacks) != 1 || sender.ackCallbacks[0] != "cb-1" {
		t.Fatalf("expected callback ack, got %#v", sender.ackCallbacks)
	}
	if len(sender.typingChats) != 1 || sender.typingChats[0] != chatID {
		t.Fatalf("expected typing for callback chat %d, got %#v", chatID, sender.typingChats)
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
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, nil, []int64{userID})

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

func TestHandler_HandleTags_WithTags(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{
		TagsList: []thingsreader.Tag{
			{Name: "work", Path: "work"},
			{Name: "home", Path: "home"},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, nil, []int64{userID})

	update := newTestUpdate(userID, chatID, "/tags")
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
	if !strings.Contains(reply.text, "Choose a tag") {
		t.Fatalf("expected choose-tag message, got %q", reply.text)
	}
	if len(reply.keyboard.InlineKeyboard) == 0 || len(reply.keyboard.InlineKeyboard[0]) == 0 {
		t.Fatalf("expected at least one button, got %#v", reply.keyboard.InlineKeyboard)
	}
	if reply.keyboard.InlineKeyboard[0][0].Text != "home" {
		t.Fatalf("expected sorted tags with home first, got %#v", reply.keyboard.InlineKeyboard)
	}
}

func TestHandler_HandleTags_Empty(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{}
	h := bot.NewHandler(sender, rec, rdr, nil, []int64{userID})

	update := newTestUpdate(userID, chatID, "/tags")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(*messages))
	}
	if !strings.Contains((*messages)[0].text, "No tags found") {
		t.Fatalf("expected empty tags message, got %q", (*messages)[0].text)
	}
}

func TestHandler_CallbackTagSelection_FirstPage(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{
		Tasks: []thingsreader.Task{
			{Title: "Task 01"}, {Title: "Task 02"},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, nil, []int64{userID})

	update := newCallbackUpdate(userID, chatID, "cb-tag-1", "tagsel:work")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sender.ackCallbacks) != 1 || sender.ackCallbacks[0] != "cb-tag-1" {
		t.Fatalf("expected callback ack, got %#v", sender.ackCallbacks)
	}
	if rdr.LastTag != "work" || rdr.LastPageOffset != 0 || rdr.LastPageLimit != 11 {
		t.Fatalf("expected tag page call for work offset=0 limit=11, got tag=%q offset=%d limit=%d", rdr.LastTag, rdr.LastPageOffset, rdr.LastPageLimit)
	}
	if len(*messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(*messages))
	}
	if (*messages)[0].withKeyboard {
		t.Fatal("expected plain text reply when there is no pagination keyboard")
	}
	if !strings.Contains((*messages)[0].text, "🏷️ work — page 1") {
		t.Fatalf("expected tag page header, got %q", (*messages)[0].text)
	}
}

func TestHandler_CallbackTagPagination_NextPage(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{
		Tasks: []thingsreader.Task{
			{Title: "Task 01"}, {Title: "Task 02"}, {Title: "Task 03"}, {Title: "Task 04"}, {Title: "Task 05"},
			{Title: "Task 06"}, {Title: "Task 07"}, {Title: "Task 08"}, {Title: "Task 09"}, {Title: "Task 10"},
			{Title: "Task 11"},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, nil, []int64{userID})

	update := newCallbackUpdate(userID, chatID, "cb-tag-2", "tagpage:work:1")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sender.ackCallbacks) != 1 || sender.ackCallbacks[0] != "cb-tag-2" {
		t.Fatalf("expected callback ack, got %#v", sender.ackCallbacks)
	}
	if rdr.LastTag != "work" || rdr.LastPageOffset != 10 || rdr.LastPageLimit != 11 {
		t.Fatalf("expected tag page call for work offset=10 limit=11, got tag=%q offset=%d limit=%d", rdr.LastTag, rdr.LastPageOffset, rdr.LastPageLimit)
	}
	if len(*messages) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(*messages))
	}
	reply := (*messages)[0]
	if !strings.Contains(reply.text, "🏷️ work — page 2") || !strings.Contains(reply.text, "11. ⬜ Task 11") {
		t.Fatalf("expected tag page 2 with task 11, got %q", reply.text)
	}
	if len(reply.keyboard.InlineKeyboard) != 1 || len(reply.keyboard.InlineKeyboard[0]) != 1 {
		t.Fatalf("expected one Prev button, got %#v", reply.keyboard.InlineKeyboard)
	}
	if reply.keyboard.InlineKeyboard[0][0].CallbackData == nil || *reply.keyboard.InlineKeyboard[0][0].CallbackData != "tagpage:work:0" {
		t.Fatalf("unexpected callback data: %#v", reply.keyboard.InlineKeyboard[0][0].CallbackData)
	}
}

func TestHandler_HandleTags_HierarchySortedByPath(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{
		TagsList: []thingsreader.Tag{
			{Name: "Urgent", Path: "Work/ClientA/Urgent"},
			{Name: "Errands", Path: "Personal/Errands"},
			{Name: "ClientA", Path: "Work/ClientA"},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, nil, []int64{userID})

	update := newTestUpdate(userID, chatID, "/tags")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reply := (*messages)[0]
	first := reply.keyboard.InlineKeyboard[0][0].Text
	second := reply.keyboard.InlineKeyboard[0][1].Text
	third := reply.keyboard.InlineKeyboard[1][0].Text
	if first != "Personal/Errands" || second != "Work/ClientA" || third != "Work/ClientA/Urgent" {
		t.Fatalf("unexpected hierarchy order: %#v", reply.keyboard.InlineKeyboard)
	}
}

func TestHandler_TaskCommand_ShowsTaskDetailsWithButtons(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	store := &RecordingTaskStore{
		tasksByNumber: map[int]thingsreader.Task{
			3: {ID: "abc123", Title: "Write report", Project: "Q1", Area: "Work", Deadline: "Friday", Tags: []string{"office"}, Completed: false},
		},
	}
	h := bot.NewHandler(sender, rec, &readertest.RecordingReader{}, store, []int64{userID})

	update := newTestUpdate(userID, chatID, "/task 3")
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
	if !strings.Contains(reply.text, "Task #3") || !strings.Contains(reply.text, "Write report") {
		t.Fatalf("unexpected task details text: %q", reply.text)
	}
}

func TestHandler_TaskCallback_Done_UpdatesViaReader(t *testing.T) {
	const authToken = "tok"
	const userID = int64(42)
	const chatID = int64(42)

	rec := &openertest.RecordingOpener{}
	sender, messages := newSenderMock()
	rdr := &readertest.RecordingReader{}
	store := &RecordingTaskStore{
		tasksByNumber: map[int]thingsreader.Task{
			5: {ID: "task-xyz", Title: "Pay rent"},
		},
	}
	h := bot.NewHandler(sender, rec, rdr, store, []int64{userID})

	update := newCallbackUpdate(userID, chatID, "cb-task-1", "taskop:done:5")
	if err := h.Handle(update); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rdr.LastUpdateID != "task-xyz" || rdr.LastCompleted == nil || !*rdr.LastCompleted {
		t.Fatalf("expected reader completion update for task-xyz, got id=%q completed=%#v", rdr.LastUpdateID, rdr.LastCompleted)
	}
	if len(*messages) != 1 || !strings.Contains((*messages)[0].text, "marked as done") {
		t.Fatalf("unexpected callback response: %#v", *messages)
	}
}
