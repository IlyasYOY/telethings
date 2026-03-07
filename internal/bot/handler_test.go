package bot_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/IlyasYOY/telethings/internal/bot"
	"github.com/IlyasYOY/telethings/internal/thingser"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestHandlerHandle_IgnoresNonCommandMessage(t *testing.T) {
	sender := NewMessageSenderMock(t)
	handler := bot.NewHandler(sender, nil, nil, nil, []int64{42})

	err := handler.Handle(tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: 42},
			Chat: &tgbotapi.Chat{ID: 1001},
			Text: "hello",
		},
	})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestHandlerHandle_IgnoresUnauthorizedCommand(t *testing.T) {
	sender := NewMessageSenderMock(t)
	handler := bot.NewHandler(sender, nil, nil, nil, []int64{42})

	err := handler.Handle(commandUpdate(777, 1001, "/start"))
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestHandlerHandle_UnknownCommand(t *testing.T) {
	sender := NewMessageSenderMock(t)
	sender.SendMock.Expect(int64(1001), "Unknown command. Use /start to see available commands.").Return(nil)

	handler := bot.NewHandler(sender, nil, nil, nil, []int64{42})
	err := handler.Handle(commandUpdate(42, 1001, "/nope"))
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestHandlerHandle_StartCommand(t *testing.T) {
	sender := NewMessageSenderMock(t)
	sender.SendMock.Set(func(chatID int64, text string) error {
		if chatID != 1001 {
			t.Fatalf("unexpected chatID: %d", chatID)
		}
		if !strings.Contains(text, "Welcome to Telethings!") {
			t.Fatalf("missing welcome text: %q", text)
		}
		if !strings.Contains(text, "/task <number>") {
			t.Fatalf("missing /task docs: %q", text)
		}
		return nil
	})

	handler := bot.NewHandler(sender, nil, nil, nil, []int64{42})
	err := handler.Handle(commandUpdate(42, 1001, "/start"))
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestHandlerHandle_AddCommandInvalidInput(t *testing.T) {
	sender := NewMessageSenderMock(t)
	sender.SendMock.Expect(int64(1001), "Usage: /add <title> [when:<value>] [deadline:<value>] [tags:<csv>] [notes:<text>]").Return(nil)

	handler := bot.NewHandler(sender, nil, nil, nil, []int64{42})
	err := handler.Handle(commandUpdate(42, 1001, "/add"))
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestHandlerHandle_TaskCommandInvalidNumber(t *testing.T) {
	sender := NewMessageSenderMock(t)
	sender.SendTypingMock.Expect(int64(1001)).Return(nil)
	sender.SendMock.Expect(int64(1001), "Usage: /task <number>").Return(nil)

	handler := bot.NewHandler(sender, nil, nil, nil, []int64{42})
	err := handler.Handle(commandUpdate(42, 1001, "/task nope"))
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestHandlerHandle_TaskCommandWithoutStore(t *testing.T) {
	sender := NewMessageSenderMock(t)
	sender.SendTypingMock.Expect(int64(1001)).Return(nil)
	sender.SendMock.Expect(int64(1001), "Task storage is not configured.").Return(nil)

	handler := bot.NewHandler(sender, nil, nil, nil, []int64{42})
	err := handler.Handle(commandUpdate(42, 1001, "/task 1"))
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestHandlerHandle_CallbackNoMessage(t *testing.T) {
	sender := NewMessageSenderMock(t)
	handler := bot.NewHandler(sender, nil, nil, nil, []int64{42})

	err := handler.Handle(tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			ID:   "cb1",
			From: &tgbotapi.User{ID: 42},
			Data: "unknown",
		},
	})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestHandlerHandle_CallbackUnauthorized(t *testing.T) {
	sender := NewMessageSenderMock(t)
	handler := bot.NewHandler(sender, nil, nil, nil, []int64{42})

	err := handler.Handle(tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			ID:   "cb1",
			From: &tgbotapi.User{ID: 777},
			Data: "unknown",
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 1001},
			},
		},
	})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestHandlerHandle_CallbackUnknownAcked(t *testing.T) {
	sender := NewMessageSenderMock(t)
	sender.AckCallbackMock.Expect("cb1").Return(nil)

	handler := bot.NewHandler(sender, nil, nil, nil, []int64{42})
	err := handler.Handle(tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			ID:   "cb1",
			From: &tgbotapi.User{ID: 42},
			Data: "unknown",
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 1001},
			},
		},
	})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestHandlerHandle_CallbackInvalidPayloads_AckOnly(t *testing.T) {
	tests := []string{
		"tagpage:%ZZ:1",
		"taskop:done:0",
		"page:anytime:-1",
	}

	for _, data := range tests {
		t.Run(data, func(t *testing.T) {
			sender := NewMessageSenderMock(t)
			sender.AckCallbackMock.Expect("cb1").Return(nil)

			handler := bot.NewHandler(sender, nil, nil, nil, []int64{42})
			err := handler.Handle(callbackUpdate(42, 1001, "cb1", data))
			if err != nil {
				t.Fatalf("Handle() error = %v", err)
			}
		})
	}
}

func TestHandlerHandle_CallbackTagSelection_PaginatedTasks(t *testing.T) {
	sender := NewMessageSenderMock(t)
	sender.AckCallbackMock.Expect("cb1").Return(nil)
	sender.SendTypingMock.Expect(int64(1001)).Return(nil)
	sender.SendWithInlineKeyboardMock.Set(func(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
		if chatID != 1001 {
			t.Fatalf("unexpected chatID: %d", chatID)
		}
		if !strings.Contains(text, "🏷️ Home/Errands — page 1") {
			t.Fatalf("unexpected text: %q", text)
		}
		if len(keyboard.InlineKeyboard) != 1 || len(keyboard.InlineKeyboard[0]) != 1 {
			t.Fatalf("unexpected keyboard rows: %#v", keyboard.InlineKeyboard)
		}
		if keyboard.InlineKeyboard[0][0].Text != "Next ➡️" {
			t.Fatalf("unexpected button text: %q", keyboard.InlineKeyboard[0][0].Text)
		}
		return nil
	})

	reader := &testThingsReader{
		tasksByTagPage: func(tag string, offset int, limit int) ([]thingser.Task, error) {
			if tag != "Home/Errands" || offset != 0 || limit != 11 {
				t.Fatalf("unexpected args: %s %d %d", tag, offset, limit)
			}
			return []thingser.Task{
				{ID: "1", Title: "One"},
				{ID: "2", Title: "Two"},
				{ID: "3", Title: "Three"},
				{ID: "4", Title: "Four"},
				{ID: "5", Title: "Five"},
				{ID: "6", Title: "Six"},
				{ID: "7", Title: "Seven"},
				{ID: "8", Title: "Eight"},
				{ID: "9", Title: "Nine"},
				{ID: "10", Title: "Ten"},
				{ID: "11", Title: "Eleven"},
			}, nil
		},
	}

	store := NewTaskStoreMock(t)
	store.SaveTaskListMock.Set(func(chatID int64, scope string, startNumber int, tasks []thingser.Task) error {
		if chatID != 1001 || scope != "tag:Home/Errands" || startNumber != 1 {
			t.Fatalf("unexpected save params: %d %s %d", chatID, scope, startNumber)
		}
		if len(tasks) != 10 {
			t.Fatalf("unexpected task count: %d", len(tasks))
		}
		return nil
	})

	handler := bot.NewHandler(sender, nil, reader, store, []int64{42})
	err := handler.Handle(callbackUpdate(42, 1001, "cb1", "tagsel:Home%2FErrands"))
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestHandlerHandle_CallbackListPagination_EmptyPageFallsBack(t *testing.T) {
	sender := NewMessageSenderMock(t)
	sender.AckCallbackMock.Expect("cb1").Return(nil)
	sender.SendTypingMock.Set(func(chatID int64) error {
		if chatID != 1001 {
			t.Fatalf("unexpected chatID: %d", chatID)
		}
		return nil
	})
	sender.SendWithInlineKeyboardMock.Set(func(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
		if chatID != 1001 {
			t.Fatalf("unexpected chatID: %d", chatID)
		}
		if !strings.Contains(text, "📋 Anytime — page 2") {
			t.Fatalf("unexpected text: %q", text)
		}
		if len(keyboard.InlineKeyboard) != 1 || len(keyboard.InlineKeyboard[0]) != 1 {
			t.Fatalf("unexpected keyboard: %#v", keyboard.InlineKeyboard)
		}
		if keyboard.InlineKeyboard[0][0].Text != "⬅️ Prev" {
			t.Fatalf("unexpected button text: %q", keyboard.InlineKeyboard[0][0].Text)
		}
		return nil
	})

	call := 0
	reader := &testThingsReader{
		tasksInListPage: func(list string, offset, limit int) ([]thingser.Task, error) {
			call++
			if list != "Anytime" || limit != 11 {
				t.Fatalf("unexpected paging args: %s %d", list, limit)
			}
			switch call {
			case 1:
				if offset != 20 {
					t.Fatalf("unexpected first offset: %d", offset)
				}
				return nil, nil
			case 2:
				if offset != 10 {
					t.Fatalf("unexpected fallback offset: %d", offset)
				}
				return []thingser.Task{{ID: "a1", Title: "Task A"}}, nil
			default:
				return nil, fmt.Errorf("unexpected call %d", call)
			}
		},
	}

	store := NewTaskStoreMock(t)
	store.SaveTaskListMock.Expect(int64(1001), "list:anytime", 11, []thingser.Task{{ID: "a1", Title: "Task A"}}).Return(nil)

	handler := bot.NewHandler(sender, nil, reader, store, []int64{42})
	err := handler.Handle(callbackUpdate(42, 1001, "cb1", "page:anytime:2"))
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestHandlerHandle_CallbackTaskOperationDone(t *testing.T) {
	sender := NewMessageSenderMock(t)
	sender.AckCallbackMock.Expect("cb1").Return(nil)
	sender.SendTypingMock.Expect(int64(1001)).Return(nil)
	sender.SendMock.Expect(int64(1001), "✅ Task marked as done").Return(nil)

	reader := &testThingsReader{
		setTaskCompleted: func(id string, completed bool) error {
			if id != "task-1" || !completed {
				t.Fatalf("unexpected completion args: %s %v", id, completed)
			}
			return nil
		},
	}
	store := NewTaskStoreMock(t)
	store.TaskByNumberMock.Expect(int64(1001), 3).Return(thingser.Task{ID: "task-1", Title: "Do thing"}, nil)

	handler := bot.NewHandler(sender, nil, reader, store, []int64{42})
	err := handler.Handle(callbackUpdate(42, 1001, "cb1", "taskop:done:3"))
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func commandUpdate(userID, chatID int64, text string) tgbotapi.Update {
	entities := []tgbotapi.MessageEntity(nil)
	if strings.HasPrefix(text, "/") {
		length := len(text)
		if space := strings.Index(text, " "); space > 0 {
			length = space
		}
		entities = []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: length},
		}
	}

	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			From:     &tgbotapi.User{ID: userID},
			Chat:     &tgbotapi.Chat{ID: chatID},
			Text:     text,
			Entities: entities,
		},
	}
}

func callbackUpdate(userID, chatID int64, callbackID, data string) tgbotapi.Update {
	return tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			ID:   callbackID,
			From: &tgbotapi.User{ID: userID},
			Data: data,
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: chatID},
			},
		},
	}
}

type testThingsReader struct {
	tasksInListPage  func(list string, offset, limit int) ([]thingser.Task, error)
	tasksByTagPage   func(tag string, offset, limit int) ([]thingser.Task, error)
	setTaskCompleted func(id string, completed bool) error
}

func (r *testThingsReader) TasksInList(string) ([]thingser.Task, error) { return nil, nil }
func (r *testThingsReader) Tags() ([]thingser.Tag, error)               { return nil, nil }
func (r *testThingsReader) AddTask(thingser.AddTaskInput) (thingser.Task, error) {
	return thingser.Task{}, nil
}
func (r *testThingsReader) SetTaskCanceled(string, bool) error { return nil }

func (r *testThingsReader) TasksInListPage(list string, offset, limit int) ([]thingser.Task, error) {
	if r.tasksInListPage != nil {
		return r.tasksInListPage(list, offset, limit)
	}
	return nil, nil
}

func (r *testThingsReader) TasksByTagPage(tag string, offset int, limit int) ([]thingser.Task, error) {
	if r.tasksByTagPage != nil {
		return r.tasksByTagPage(tag, offset, limit)
	}
	return nil, nil
}

func (r *testThingsReader) SetTaskCompleted(id string, completed bool) error {
	if r.setTaskCompleted != nil {
		return r.setTaskCompleted(id, completed)
	}
	return nil
}
