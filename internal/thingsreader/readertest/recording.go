package readertest

import "github.com/IlyasYOY/telethings/internal/thingsreader"

// RecordingReader is a test double for the thingsReader interface.
// Set Tasks to control the returned task list; set Err to simulate errors.
type RecordingReader struct {
	Tasks          []thingsreader.Task
	Err            error
	LastPageList   string
	LastPageOffset int
	LastPageLimit  int
}

// TasksInList returns the pre-configured Tasks and Err, ignoring list.
func (r *RecordingReader) TasksInList(_ string) ([]thingsreader.Task, error) {
	return r.Tasks, r.Err
}

// TasksInListPage returns a page from pre-configured Tasks and records pagination args.
func (r *RecordingReader) TasksInListPage(list string, offset, limit int) ([]thingsreader.Task, error) {
	r.LastPageList = list
	r.LastPageOffset = offset
	r.LastPageLimit = limit
	if r.Err != nil {
		return nil, r.Err
	}
	if limit <= 0 || offset >= len(r.Tasks) {
		return nil, nil
	}
	if offset < 0 {
		offset = 0
	}
	end := offset + limit
	if end > len(r.Tasks) {
		end = len(r.Tasks)
	}
	return append([]thingsreader.Task(nil), r.Tasks[offset:end]...), nil
}
