package bot

import "github.com/IlyasYOY/telethings/internal/thingsreader"

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock@v3.4.7 -i thingsReader -o reader_mock_test.go -g

// thingsReader reads task lists from Things 3.
type thingsReader interface {
	TasksInList(list string) ([]thingsreader.Task, error)
	TasksInListPage(list string, offset, limit int) ([]thingsreader.Task, error)
	TasksByTagPage(tag string, offset, limit int) ([]thingsreader.Task, error)
	Tags() ([]thingsreader.Tag, error)
	AddTask(input thingsreader.AddTaskInput) (thingsreader.Task, error)
	SetTaskCompleted(id string, completed bool) error
	SetTaskCanceled(id string, canceled bool) error
}
