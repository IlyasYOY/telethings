-- +goose Up
CREATE TABLE IF NOT EXISTS task_list_state (
  chat_id INTEGER PRIMARY KEY,
  scope TEXT NOT NULL,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS task_list_items (
  chat_id INTEGER NOT NULL,
  item_number INTEGER NOT NULL,
  task_id TEXT NOT NULL,
  title TEXT NOT NULL,
  project TEXT NOT NULL,
  area TEXT NOT NULL,
  deadline TEXT NOT NULL,
  tags_csv TEXT NOT NULL,
  completed INTEGER NOT NULL,
  PRIMARY KEY (chat_id, item_number)
);

CREATE INDEX IF NOT EXISTS idx_task_list_items_chat ON task_list_items(chat_id);

-- +goose Down
DROP INDEX IF EXISTS idx_task_list_items_chat;
DROP TABLE IF EXISTS task_list_items;
DROP TABLE IF EXISTS task_list_state;
