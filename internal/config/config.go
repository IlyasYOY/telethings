// Package config provides configuration loading from environment variables.
package config

import (
	"errors"
	"os"
)

const (
	envTelegramToken   = "TELETHINGS_TELEGRAM_TOKEN"
	envThingsAuthToken = "TELETHINGS_THINGS_AUTH_TOKEN"
)

var (
	ErrMissingTelegramToken   = errors.New("telegram token not set: " + envTelegramToken)
	ErrMissingThingsAuthToken = errors.New("things auth token not set: " + envThingsAuthToken)
)

// Config holds the runtime configuration for the bot.
type Config struct {
	TelegramToken   string
	ThingsAuthToken string
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

	return &Config{
		TelegramToken:   telegramToken,
		ThingsAuthToken: thingsAuthToken,
	}, nil
}
