package db

import (
	"errors"
	"testing"

	"github.com/IlyasYOY/telethings/internal/thingsreader"
)

func TestOpenAndMigrate_DefaultDSNSuccess(t *testing.T) {
	conn, err := OpenAndMigrate("")
	if err != nil {
		t.Fatalf("OpenAndMigrate returned error: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	for _, table := range []string{"task_list_state", "task_list_items"} {
		var name string
		if err := conn.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&name); err != nil {
			t.Fatalf("table %q not found: %v", table, err)
		}
	}
}

func TestOpenAndMigrate_Failure(t *testing.T) {
	conn, err := OpenAndMigrate("file:telethings-invalid?mode=invalid")
	if err == nil {
		_ = conn.Close()
		t.Fatal("expected error, got nil")
	}
}

func TestTaskStore_SaveTaskListAndTaskByNumber(t *testing.T) {
	conn, err := OpenAndMigrate("file:taskstore-save-task-by-number?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("OpenAndMigrate returned error: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	store := NewTaskStore(conn)
	tasks := []thingsreader.Task{
		{ID: "1", Title: "first", Project: "p1", Area: "a1", Deadline: "2026-01-01", Tags: []string{"work", "urgent"}, Completed: true},
		{ID: "2", Title: "second", Project: "p2", Area: "a2", Deadline: "", Tags: nil, Completed: false},
	}
	if err := store.SaveTaskList(42, "today", 10, tasks); err != nil {
		t.Fatalf("SaveTaskList returned error: %v", err)
	}

	got, err := store.TaskByNumber(42, 10)
	if err != nil {
		t.Fatalf("TaskByNumber returned error: %v", err)
	}
	if got.ID != "1" || got.Title != "first" || got.Project != "p1" || got.Area != "a1" || got.Deadline != "2026-01-01" {
		t.Fatalf("unexpected task: %+v", got)
	}
	if len(got.Tags) != 2 || got.Tags[0] != "work" || got.Tags[1] != "urgent" {
		t.Fatalf("unexpected tags: %#v", got.Tags)
	}
	if !got.Completed {
		t.Fatal("expected completed task")
	}

	got, err = store.TaskByNumber(42, 11)
	if err != nil {
		t.Fatalf("TaskByNumber returned error: %v", err)
	}
	if len(got.Tags) != 0 {
		t.Fatalf("expected no tags, got: %#v", got.Tags)
	}
	if got.Completed {
		t.Fatal("expected not completed task")
	}
}

func TestTaskStore_TaskByNumberNotFound(t *testing.T) {
	conn, err := OpenAndMigrate("file:taskstore-not-found?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("OpenAndMigrate returned error: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	store := NewTaskStore(conn)
	_, err = store.TaskByNumber(100, 1)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}
