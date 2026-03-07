// Package thingsreader provides implementations for reading data from Things 3.
package thingser

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Task represents a Things 3 task with optional metadata.
type Task struct {
	ID        string
	Title     string
	Project   string
	Deadline  string
	Tags      []string
	Area      string
	Completed bool
	Canceled  bool
}

// Tag represents a Things 3 tag with full hierarchy path.
type Tag struct {
	Name string
	Path string
}

// AppleScriptReader reads task lists from Things 3 via AppleScript.
type AppleScriptReader struct{}

// Tags returns all tags from Things 3 with hierarchy paths.
func (AppleScriptReader) Tags() ([]Tag, error) {
	script := `tell application "Things3"
set fieldSep to ASCII character 31
set lineSep to ASCII character 10
set rows to {}
repeat with tg in tags
	set tagName to name of tg
	set pathParts to {tagName}
	set currentTag to tg
	repeat
		try
			set parentTag to parent tag of currentTag
			if parentTag is missing value then exit repeat
			set beginning of pathParts to name of parentTag
			set currentTag to parentTag
		on error
			exit repeat
		end try
	end repeat

	set AppleScript's text item delimiters to "/"
	set tagPath to pathParts as string
	set AppleScript's text item delimiters to fieldSep
	set end of rows to tagName & fieldSep & tagPath
end repeat

set AppleScript's text item delimiters to lineSep
return rows as string
end tell`
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return nil, err
	}
	return parseTagsListOutput(out), nil
}

// TasksInList returns all to-dos in the named Things 3 list with metadata.
func (AppleScriptReader) TasksInList(list string) ([]Task, error) {
	script := `tell application "Things3"
set fieldSep to ASCII character 31
set lineSep to ASCII character 10
set rows to {}
repeat with t in to dos of list "` + list + `"
	set taskID to id of t as string
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

	set statusText to ""
	try
		set taskStatus to status of t as string
		set statusText to taskStatus
	end try

	set AppleScript's text item delimiters to fieldSep
	set end of rows to taskID & fieldSep & taskName & fieldSep & projectName & fieldSep & deadlineText & fieldSep & tagsText & fieldSep & areaName & fieldSep & statusText
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
	set taskID to id of t as string
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

	set statusText to ""
	try
		set taskStatus to status of t as string
		set statusText to taskStatus
	end try

	set AppleScript's text item delimiters to fieldSep
	set end of rows to taskID & fieldSep & taskName & fieldSep & projectName & fieldSep & deadlineText & fieldSep & tagsText & fieldSep & areaName & fieldSep & statusText
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

// TasksByTagPage returns up to limit to-dos with the given tag path, starting at offset.
func (AppleScriptReader) TasksByTagPage(tagPath string, offset, limit int) ([]Task, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		return nil, nil
	}

	script := tasksByTagPageScript(tagPath, offset, limit)
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return nil, err
	}
	return parseTasksOutput(out), nil
}

func tasksByTagPageScript(tagPath string, offset, limit int) string {
	escapedTagPath := escapeAppleScriptString(tagPath)
	escapedTagPathPrefix := escapeAppleScriptString(tagPath + "/")
	return fmt.Sprintf(`tell application "Things3"
set fieldSep to ASCII character 31
set lineSep to ASCII character 10
set rows to {}
set selectedTags to {}
set seenIDs to {}
set selectedCount to 0
set neededCount to %d + %d
set doneCollecting to false

repeat with tg in tags
	set pathParts to {name of tg}
	set currentTag to tg
	repeat
		try
			set parentTag to parent tag of currentTag
			if parentTag is missing value then exit repeat
			set beginning of pathParts to name of parentTag
			set currentTag to parentTag
		on error
			exit repeat
		end try
	end repeat
	set AppleScript's text item delimiters to "/"
	set candidatePath to pathParts as string
	if candidatePath is "%s" or candidatePath starts with "%s" then
		set end of selectedTags to tg
	end if
end repeat

if neededCount <= 0 then
	return ""
end if

repeat with tg in selectedTags
	try
		repeat with t in to dos of tg
			set todoID to id of t as string
			if seenIDs does not contain todoID then
				set end of seenIDs to todoID
				set selectedCount to selectedCount + 1
				if selectedCount > %d and selectedCount <= neededCount then
					set taskID to id of t as string
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

					set statusText to ""
					try
						set taskStatus to status of t as string
						set statusText to taskStatus
					end try

					set AppleScript's text item delimiters to fieldSep
					set end of rows to taskID & fieldSep & taskName & fieldSep & projectName & fieldSep & deadlineText & fieldSep & tagsText & fieldSep & areaName & fieldSep & statusText
				end if

				if selectedCount >= neededCount then
					set doneCollecting to true
					exit repeat
				end if
			end if
		end repeat
	on error
		-- ignore invalid tag entries and continue
	end try
	if doneCollecting then exit repeat
end repeat

set AppleScript's text item delimiters to lineSep
return rows as string
end tell`, offset, limit, escapedTagPath, escapedTagPathPrefix, offset)
}

func parseTasksOutput(out []byte) []Task {
	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return nil
	}

	lines := strings.Split(raw, "\n")
	tasks := make([]Task, 0, len(lines))
	for _, line := range lines {
		fields := strings.SplitN(line, string(rune(31)), 7)
		if len(fields) != 7 {
			continue
		}
		tasks = append(tasks, Task{
			ID:        fields[0],
			Title:     fields[1],
			Project:   fields[2],
			Deadline:  fields[3],
			Tags:      parseTags(fields[4]),
			Area:      fields[5],
			Completed: parseCompleted(fields[6]),
			Canceled:  parseCanceled(fields[6]),
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

func parseCanceled(raw string) bool {
	value := strings.ToLower(strings.TrimSpace(raw))
	return value == "canceled" || value == "cancelled"
}

func parseTagsListOutput(out []byte) []Tag {
	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return nil
	}
	lines := strings.Split(raw, "\n")
	tags := make([]Tag, 0, len(lines))
	for _, line := range lines {
		fields := strings.SplitN(line, string(rune(31)), 2)
		if len(fields) != 2 {
			continue
		}
		name := strings.TrimSpace(fields[0])
		path := strings.TrimSpace(fields[1])
		if name == "" || path == "" {
			continue
		}
		tags = append(tags, Tag{Name: name, Path: path})
	}
	return tags
}

func escapeAppleScriptString(value string) string {
	s := strings.ReplaceAll(value, `\`, `\\`)
	return strings.ReplaceAll(s, `"`, `\"`)
}
