package bot

import "github.com/IlyasYOY/telethings/internal/thingsreader"

type taskStore interface {
	SaveTaskList(chatID int64, scope string, startNumber int, tasks []thingsreader.Task) error
	TaskByNumber(chatID int64, number int) (thingsreader.Task, error)
}
