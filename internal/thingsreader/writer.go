package thingsreader

import (
	"fmt"
	"os/exec"
	"strings"
)

type AddTaskInput struct {
	Title    string
	When     string
	Deadline string
	Tags     []string
	Notes    string
}

func (AppleScriptReader) AddTask(input AddTaskInput) (Task, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return Task{}, fmt.Errorf("empty title")
	}

	script := buildAddTaskScript(input)
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return Task{}, err
	}
	id := strings.TrimSpace(string(out))
	if id == "" {
		return Task{}, fmt.Errorf("empty task id from add script")
	}
	return Task{
		ID:       id,
		Title:    title,
		Deadline: input.Deadline,
		Tags:     append([]string(nil), input.Tags...),
	}, nil
}

func (AppleScriptReader) SetTaskCompleted(id string, completed bool) error {
	status := "open"
	if completed {
		status = "completed"
	}
	script := fmt.Sprintf(`tell application "Things3"
set t to to do id "%s"
set status of t to %s
end tell`, escapeAppleScriptString(id), status)
	return exec.Command("osascript", "-e", script).Run()
}

func (AppleScriptReader) SetTaskCanceled(id string, canceled bool) error {
	status := "open"
	if canceled {
		status = "canceled"
	}
	script := fmt.Sprintf(`tell application "Things3"
set t to to do id "%s"
set status of t to %s
end tell`, escapeAppleScriptString(id), status)
	return exec.Command("osascript", "-e", script).Run()
}

func buildAddTaskScript(input AddTaskInput) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "tell application \"Things3\"\n")
	fmt.Fprintf(&sb, "set newTodo to make new to do with properties {name:\"%s\"}\n", escapeAppleScriptString(strings.TrimSpace(input.Title)))

	if notes := strings.TrimSpace(input.Notes); notes != "" {
		fmt.Fprintf(&sb, "set notes of newTodo to \"%s\"\n", escapeAppleScriptString(notes))
	}
	if deadline := strings.TrimSpace(input.Deadline); deadline != "" {
		fmt.Fprintf(&sb, "try\n")
		fmt.Fprintf(&sb, "set due date of newTodo to date \"%s\"\n", escapeAppleScriptString(deadline))
		fmt.Fprintf(&sb, "end try\n")
	}
	if when := strings.ToLower(strings.TrimSpace(input.When)); when != "" {
		switch when {
		case "today":
			fmt.Fprintf(&sb, "try\n")
			fmt.Fprintf(&sb, "move newTodo to list \"Today\"\n")
			fmt.Fprintf(&sb, "on error\n")
			fmt.Fprintf(&sb, "try\n")
			fmt.Fprintf(&sb, "set activation date of newTodo to (current date)\n")
			fmt.Fprintf(&sb, "end try\n")
			fmt.Fprintf(&sb, "end try\n")
		case "tomorrow":
			fmt.Fprintf(&sb, "try\n")
			fmt.Fprintf(&sb, "move newTodo to list \"Tomorrow\"\n")
			fmt.Fprintf(&sb, "on error\n")
			fmt.Fprintf(&sb, "try\n")
			fmt.Fprintf(&sb, "set activation date of newTodo to ((current date) + 1 * days)\n")
			fmt.Fprintf(&sb, "end try\n")
			fmt.Fprintf(&sb, "end try\n")
		}
	}
	if len(input.Tags) > 0 {
		fmt.Fprintf(&sb, "set tag names of newTodo to {%s}\n", appleScriptQuotedList(input.Tags))
	}

	fmt.Fprintf(&sb, "delay %.1f\n", 0.2)
	fmt.Fprintf(&sb, "return (id of newTodo as string)\n")
	fmt.Fprintf(&sb, "end tell\n")
	return sb.String()
}

func appleScriptQuotedList(items []string) string {
	quoted := make([]string, 0, len(items))
	for _, item := range items {
		v := strings.TrimSpace(item)
		if v == "" {
			continue
		}
		quoted = append(quoted, `"`+escapeAppleScriptString(v)+`"`)
	}
	return strings.Join(quoted, ", ")
}
