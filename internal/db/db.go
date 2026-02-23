package db

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

const defaultDSN = "file:telethings?mode=memory&cache=shared"

//go:embed migrations/*.sql
var migrationsFS embed.FS

func OpenAndMigrate(dsn string) (*sql.DB, error) {
	if dsn == "" {
		dsn = defaultDSN
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
