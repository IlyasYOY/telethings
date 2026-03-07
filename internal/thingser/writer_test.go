package thingser

import (
	"strings"
	"testing"
)

func TestBuildAddTaskScript_WhenTodayWrappedInTry(t *testing.T) {
	script := buildAddTaskScript(AddTaskInput{Title: "Task", When: "today"})
	if !strings.Contains(script, `move newTodo to list "Today"`) {
		t.Fatalf("expected Today move in script, got:\n%s", script)
	}
	if !strings.Contains(script, "on error") || !strings.Contains(script, "set activation date of newTodo to (current date)") {
		t.Fatalf("expected fallback activation date in script, got:\n%s", script)
	}
}

func TestBuildAddTaskScript_DeadlineWrappedInTry(t *testing.T) {
	script := buildAddTaskScript(AddTaskInput{Title: "Task", Deadline: "2026-12-31"})
	if !strings.Contains(script, "try\nset due date of newTodo to date") {
		t.Fatalf("expected due date assignment wrapped in try, got:\n%s", script)
	}
}
