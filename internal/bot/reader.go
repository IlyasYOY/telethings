package bot

// thingsReader reads task lists from Things 3.
type thingsReader interface {
	TasksInList(list string) ([]string, error)
}
