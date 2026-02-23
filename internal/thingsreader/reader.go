// Package thingsreader provides implementations for reading data from Things 3.
package thingsreader

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Task represents a Things 3 task with optional metadata.
type Task struct {
	Title     string
	Project   string
	Deadline  string
	Tags      []string
	Area      string
	Completed bool
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
		set d to due date of t
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

	set completedText to "false"
	try
		set taskStatus to status of t as string
		if taskStatus is "completed" then set completedText to "true"
	end try

	set AppleScript's text item delimiters to fieldSep
	set end of rows to taskName & fieldSep & projectName & fieldSep & deadlineText & fieldSep & tagsText & fieldSep & areaName & fieldSep & completedText
end repeat

set AppleScript's text item delimiters to lineSep
return rows as string
end tell`
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return nil, err
	}
	return parseTasksOutput(out), nil
}

// TasksInListPage returns up to limit to-dos from the named Things 3 list, starting at offset.
func (AppleScriptReader) TasksInListPage(list string, offset, limit int) ([]Task, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		return nil, nil
	}

	script := fmt.Sprintf(`tell application "Things3"
set fieldSep to ASCII character 31
set lineSep to ASCII character 10
set rows to {}
set allTodos to to dos of list "%s"
set totalCount to count of allTodos
set startIndex to %d + 1
set endIndex to %d + %d

if startIndex > totalCount then
	return ""
end if
if endIndex > totalCount then
	set endIndex to totalCount
end if

repeat with t in items startIndex thru endIndex of allTodos
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
		set d to due date of t
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

	set completedText to "false"
	try
		set taskStatus to status of t as string
		if taskStatus is "completed" then set completedText to "true"
	end try

	set AppleScript's text item delimiters to fieldSep
	set end of rows to taskName & fieldSep & projectName & fieldSep & deadlineText & fieldSep & tagsText & fieldSep & areaName & fieldSep & completedText
end repeat

set AppleScript's text item delimiters to lineSep
return rows as string
end tell`, list, offset, offset, limit)

	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return nil, err
	}
	return parseTasksOutput(out), nil
}

func parseTasksOutput(out []byte) []Task {
	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return nil
	}

	lines := strings.Split(raw, "\n")
	tasks := make([]Task, 0, len(lines))
	for _, line := range lines {
		fields := strings.SplitN(line, string(rune(31)), 6)
		if len(fields) != 6 {
			continue
		}
		tasks = append(tasks, Task{
			Title:     fields[0],
			Project:   fields[1],
			Deadline:  fields[2],
			Tags:      parseTags(fields[3]),
			Area:      fields[4],
			Completed: parseCompleted(fields[5]),
		})
	}
	return tasks
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

func parseCompleted(raw string) bool {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "true" || value == "completed" || value == "yes" {
		return true
	}
	n, err := strconv.ParseBool(value)
	return err == nil && n
}
