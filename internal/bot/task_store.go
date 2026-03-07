package bot

import "github.com/IlyasYOY/telethings/internal/thingser"

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock@v3.4.7 -i taskStore -o task_store_mock_test.go -p bot_test -g

type taskStore interface {
	SaveTaskList(chatID int64, scope string, startNumber int, tasks []thingser.Task) error
	TaskByNumber(chatID int64, number int) (thingser.Task, error)
}
