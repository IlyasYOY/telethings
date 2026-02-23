package bot

import "github.com/IlyasYOY/telethings/internal/thingsreader"

// thingsReader reads task lists from Things 3.
type thingsReader interface {
	TasksInList(list string) ([]thingsreader.Task, error)
	TasksInListPage(list string, offset, limit int) ([]thingsreader.Task, error)
}
