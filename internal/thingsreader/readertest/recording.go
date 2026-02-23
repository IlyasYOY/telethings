package readertest

import "github.com/IlyasYOY/telethings/internal/thingsreader"

// RecordingReader is a test double for the thingsReader interface.
// Set Tasks to control the returned task list; set Err to simulate errors.
type RecordingReader struct {
	Tasks          []thingsreader.Task
	TagsList       []thingsreader.Tag
	Err            error
	AddErr         error
	UpdateErr      error
	LastPageList   string
	LastPageOffset int
	LastPageLimit  int
	LastTag        string
	LastAddInput   thingsreader.AddTaskInput
	LastUpdateID   string
	LastCompleted  *bool
	LastCanceled   *bool
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

// TasksByTagPage returns a page from pre-configured Tasks and records tag pagination args.
func (r *RecordingReader) TasksByTagPage(tag string, offset, limit int) ([]thingsreader.Task, error) {
	r.LastTag = tag
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

// Tags returns the pre-configured tag list and Err.
func (r *RecordingReader) Tags() ([]thingsreader.Tag, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	return append([]thingsreader.Tag(nil), r.TagsList...), nil
}

// AddTask records input and returns a task with configured ID.
func (r *RecordingReader) AddTask(input thingsreader.AddTaskInput) (thingsreader.Task, error) {
	r.LastAddInput = input
	if r.AddErr != nil {
		return thingsreader.Task{}, r.AddErr
	}
	id := "new-task-id"
	if len(r.Tasks) > 0 && r.Tasks[0].ID != "" {
		id = r.Tasks[0].ID
	}
	task := thingsreader.Task{
		ID:       id,
		Title:    input.Title,
		Deadline: input.Deadline,
		Tags:     append([]string(nil), input.Tags...),
	}
	return task, nil
}

func (r *RecordingReader) SetTaskCompleted(id string, completed bool) error {
	r.LastUpdateID = id
	r.LastCompleted = &completed
	if r.UpdateErr != nil {
		return r.UpdateErr
	}
	return nil
}

func (r *RecordingReader) SetTaskCanceled(id string, canceled bool) error {
	r.LastUpdateID = id
	r.LastCanceled = &canceled
	if r.UpdateErr != nil {
		return r.UpdateErr
	}
	return nil
}
