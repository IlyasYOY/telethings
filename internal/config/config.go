// Package config provides configuration loading from environment variables.
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

const (
	envTelegramToken   = "TELETHINGS_TELEGRAM_TOKEN"
	envThingsAuthToken = "TELETHINGS_THINGS_AUTH_TOKEN"
	envAllowedUserIDs  = "TELETHINGS_ALLOWED_USER_IDS"
	envDBDSN           = "TELETHINGS_DB_DSN"
)

var (
	ErrMissingTelegramToken   = errors.New("telegram token not set: " + envTelegramToken)
	ErrMissingThingsAuthToken = errors.New("things auth token not set: " + envThingsAuthToken)
	ErrMissingAllowedUserIDs  = errors.New("allowed user IDs not set: " + envAllowedUserIDs)
	ErrInvalidAllowedUserIDs  = errors.New("invalid user IDs: must be comma-separated integers")
)

// Config holds the runtime configuration for the bot.
type Config struct {
	TelegramToken   string
	ThingsAuthToken string
	AllowedUserIDs  []int64
	DBDSN           string
}

// FromEnv reads configuration from environment variables.
// It returns an error if any required variable is missing.
func FromEnv() (*Config, error) {
	telegramToken := os.Getenv(envTelegramToken)
	if telegramToken == "" {
		return nil, ErrMissingTelegramToken
	}

	thingsAuthToken := os.Getenv(envThingsAuthToken)
	if thingsAuthToken == "" {
		return nil, ErrMissingThingsAuthToken
	}

	allowedUserIDsStr := os.Getenv(envAllowedUserIDs)
	if allowedUserIDsStr == "" {
		return nil, ErrMissingAllowedUserIDs
	}

	allowedUserIDs, err := parseUserIDs(allowedUserIDsStr)
	if err != nil {
		return nil, ErrInvalidAllowedUserIDs
	}

	return &Config{
		TelegramToken:   telegramToken,
		ThingsAuthToken: thingsAuthToken,
		AllowedUserIDs:  allowedUserIDs,
		DBDSN:           strings.TrimSpace(os.Getenv(envDBDSN)),
	}, nil
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
