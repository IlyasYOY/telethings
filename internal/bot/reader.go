package bot

import "github.com/IlyasYOY/telethings/internal/thingser"

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock@v3.4.7 -i thingsReader -o reader_mock_test.go -g

// thingsReader reads task lists from Things 3.
type thingsReader interface {
	TasksInList(list string) ([]thingser.Task, error)
	TasksInListPage(list string, offset, limit int) ([]thingser.Task, error)
	TasksByTagPage(tag string, offset, limit int) ([]thingser.Task, error)
	Tags() ([]thingser.Tag, error)
	AddTask(input thingser.AddTaskInput) (thingser.Task, error)
	SetTaskCompleted(id string, completed bool) error
	SetTaskCanceled(id string, canceled bool) error
}
