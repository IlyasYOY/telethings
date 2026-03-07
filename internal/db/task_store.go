package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/IlyasYOY/telethings/internal/thingser"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskStore struct {
	db *sql.DB
}

func NewTaskStore(db *sql.DB) *TaskStore {
	return &TaskStore{db: db}
}

func (s *TaskStore) SaveTaskList(chatID int64, scope string, startNumber int, tasks []thingser.Task) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`INSERT INTO task_list_state(chat_id, scope, updated_at)
VALUES(?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(chat_id) DO UPDATE SET scope=excluded.scope, updated_at=CURRENT_TIMESTAMP`, chatID, scope); err != nil {
		return fmt.Errorf("upsert list state: %w", err)
	}

	if _, err := tx.Exec(`DELETE FROM task_list_items WHERE chat_id = ?`, chatID); err != nil {
		return fmt.Errorf("clear list items: %w", err)
	}

	stmt, err := tx.Prepare(`INSERT INTO task_list_items(chat_id, item_number, task_id, title, project, area, deadline, tags_csv, completed)
VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare insert list item: %w", err)
	}
	defer stmt.Close()

	for i, task := range tasks {
		itemNumber := startNumber + i
		if _, err := stmt.Exec(chatID, itemNumber, task.ID, task.Title, task.Project, task.Area, task.Deadline, strings.Join(task.Tags, ","), boolToInt(task.Completed)); err != nil {
			return fmt.Errorf("insert list item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (s *TaskStore) TaskByNumber(chatID int64, number int) (thingser.Task, error) {
	row := s.db.QueryRow(`SELECT task_id, title, project, area, deadline, tags_csv, completed
FROM task_list_items
WHERE chat_id = ? AND item_number = ?`, chatID, number)

	var task thingser.Task
	var tagsCSV string
	var completed int
	if err := row.Scan(&task.ID, &task.Title, &task.Project, &task.Area, &task.Deadline, &tagsCSV, &completed); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return thingser.Task{}, ErrTaskNotFound
		}
		return thingser.Task{}, fmt.Errorf("select list item: %w", err)
	}
	if tagsCSV != "" {
		task.Tags = strings.Split(tagsCSV, ",")
	}
	task.Completed = completed == 1
	return task, nil
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
