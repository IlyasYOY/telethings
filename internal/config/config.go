// Package config provides configuration loading from environment variables.
package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	envTelegramToken  = "TELETHINGS_TELEGRAM_TOKEN"
	envAllowedUserIDs = "TELETHINGS_ALLOWED_USER_IDS"
	envDBDSN          = "TELETHINGS_DB_DSN"
)

var (
	ErrMissingTelegramToken  = errors.New("telegram token not set: " + envTelegramToken)
	ErrMissingAllowedUserIDs = errors.New("allowed user IDs not set: " + envAllowedUserIDs)
	ErrInvalidAllowedUserIDs = errors.New("invalid user IDs: must be comma-separated integers")
)

// Config holds the runtime configuration for the bot.
type Config struct {
	TelegramToken  string
	AllowedUserIDs []int64
	DBDSN          string
}

// FromEnv reads configuration from environment variables.
// It returns an error if any required variable is missing.
func FromEnv() (*Config, error) {
	telegramToken := os.Getenv(envTelegramToken)
	if telegramToken == "" {
		return nil, ErrMissingTelegramToken
	}

	allowedUserIDsStr := os.Getenv(envAllowedUserIDs)
	if allowedUserIDsStr == "" {
		return nil, ErrMissingAllowedUserIDs
	}

	allowedUserIDs, err := parseUserIDs(allowedUserIDsStr)
	if err != nil {
		return nil, ErrInvalidAllowedUserIDs
	}

	dsn := strings.TrimSpace(os.Getenv(envDBDSN))
	if dsn == "" {
		dsn = defaultDBDSN()
	}

	return &Config{
		TelegramToken:  telegramToken,
		AllowedUserIDs: allowedUserIDs,
		DBDSN:          dsn,
	}, nil
}

// defaultDBDSN returns the default SQLite DSN using the XDG data directory.
// It falls back to in-memory SQLite if the home directory cannot be determined.
func defaultDBDSN() string {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "file:telethings?mode=memory&cache=shared"
		}
		dataHome = filepath.Join(home, ".local", "share")
	}
	return "file:" + filepath.Join(dataHome, "telethings", "telethings.db")
}

// parseUserIDs converts a comma-separated string of user IDs to a slice of int64.
func parseUserIDs(s string) ([]int64, error) {
	parts := strings.Split(strings.TrimSpace(s), ",")
	ids := make([]int64, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
