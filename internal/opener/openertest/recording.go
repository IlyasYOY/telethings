package openertest

// RecordingOpener records opened URLs without invoking any OS command.
// Useful in tests.
type RecordingOpener struct {
	URLs []string
}

// Open records the URL.
func (r *RecordingOpener) Open(url string) error {
	r.URLs = append(r.URLs, url)
	return nil
}
