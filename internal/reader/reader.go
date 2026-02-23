// Package reader provides implementations for reading data from Things 3.
package reader

import (
	"os/exec"
	"strings"
)

// Task represents a Things 3 task with optional metadata.
type Task struct {
	Title    string
	Project  string
	Deadline string
	Tags     []string
	Area     string
}

// AppleScriptReader reads task lists from Things 3 via AppleScript.
type AppleScriptReader struct{}

// TasksInList returns all to-dos in the named Things 3 list with metadata.
func (AppleScriptReader) TasksInList(list string) ([]Task, error) {
	script := `tell application "Things3"
set fieldSep to ASCII character 31
set lineSep to ASCII character 10
set rows to {}
repeat with t in to dos of list "` + list + `"
	set taskName to name of t

	set projectName to ""
	try
		set p to project of t
		if p is not missing value then set projectName to name of p
	end try

	set areaName to ""
	try
		set a to area of t
		if a is not missing value then set areaName to name of a
	end try

	set deadlineText to ""
	try
		set d to deadline of t
		if d is not missing value then set deadlineText to (d as string)
	end try

	set tagsText to ""
	try
		set tagList to tag names of t
		if (count of tagList) > 0 then
			set AppleScript's text item delimiters to ","
			set tagsText to tagList as string
		end if
	end try

	set AppleScript's text item delimiters to fieldSep
	set end of rows to taskName & fieldSep & projectName & fieldSep & deadlineText & fieldSep & tagsText & fieldSep & areaName
end repeat

set AppleScript's text item delimiters to lineSep
return rows as string
end tell`
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return nil, err
	}
	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return nil, nil
	}

	lines := strings.Split(raw, "\n")
	tasks := make([]Task, 0, len(lines))
	for _, line := range lines {
		fields := strings.SplitN(line, string(rune(31)), 5)
		if len(fields) != 5 {
			continue
		}
		tasks = append(tasks, Task{
			Title:    fields[0],
			Project:  fields[1],
			Deadline: fields[2],
			Tags:     parseTags(fields[3]),
			Area:     fields[4],
		})
	}
	return tasks, nil
}

func parseTags(raw string) []string {
	if raw == "" {
		return nil
	}
	items := strings.Split(raw, ",")
	tags := make([]string, 0, len(items))
	for _, item := range items {
		tag := strings.TrimSpace(item)
		if tag == "" {
			continue
		}
		tags = append(tags, tag)
	}
	return tags
}
