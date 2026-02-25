package bot

import "github.com/IlyasYOY/telethings/internal/thingsreader"

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock@v3.4.7 -i taskStore -o task_store_mock_test.go -p bot_test -g

type taskStore interface {
	SaveTaskList(chatID int64, scope string, startNumber int, tasks []thingsreader.Task) error
	TaskByNumber(chatID int64, number int) (thingsreader.Task, error)
}
