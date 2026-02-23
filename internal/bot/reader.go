package bot

import "github.com/IlyasYOY/telethings/internal/reader"

// thingsReader reads task lists from Things 3.
type thingsReader interface {
	TasksInList(list string) ([]reader.Task, error)
}
