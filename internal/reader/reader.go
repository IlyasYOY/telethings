// Package reader provides implementations for reading data from Things 3.
package reader

import (
	"os/exec"
	"strings"
)

// AppleScriptReader reads task lists from Things 3 via AppleScript.
type AppleScriptReader struct{}

// TasksInList returns the names of all to-dos in the named Things 3 list.
func (AppleScriptReader) TasksInList(list string) ([]string, error) {
	script := `tell application "Things3"
set taskList to name of to dos of list "` + list + `"
set AppleScript's text item delimiters to ASCII character 10
return taskList as string
end tell`
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return nil, err
	}
	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return nil, nil
	}
	parts := strings.Split(raw, "\n")
	return parts, nil
}
