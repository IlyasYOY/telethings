package readertest

// RecordingReader is a test double for the thingsReader interface.
// Set Tasks to control the returned task list; set Err to simulate errors.
type RecordingReader struct {
	Tasks []string
	Err   error
}

// TasksInList returns the pre-configured Tasks and Err, ignoring list.
func (r *RecordingReader) TasksInList(_ string) ([]string, error) {
	return r.Tasks, r.Err
}
