package db

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func OpenAndMigrate(dsn string) (*sql.DB, error) {
	if err := ensureDBDir(dsn); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("set goose dialect: %w", err)
	}
	if err := goose.Up(conn, "migrations"); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}
	return conn, nil
}

// ensureDBDir creates the parent directory for file-based SQLite DSNs.
// It is a no-op for in-memory DSNs.
func ensureDBDir(dsn string) error {
	if !strings.HasPrefix(dsn, "file:") {
		return nil
	}
	if strings.Contains(dsn, "mode=memory") {
		return nil
	}
	// Strip "file:" prefix and any query parameters to get the file path.
	path := strings.TrimPrefix(dsn, "file:")
	if idx := strings.Index(path, "?"); idx >= 0 {
		path = path[:idx]
	}
	if path == "" {
		return nil
	}
	return os.MkdirAll(filepath.Dir(path), 0700)
}
