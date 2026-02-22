package config_test

import (
	"errors"
	"testing"

	"github.com/IlyasYOY/telethings/internal/config"
)

func TestFromEnv_Success(t *testing.T) {
	t.Setenv("TELETHINGS_TELEGRAM_TOKEN", "tok123")
	t.Setenv("TELETHINGS_THINGS_AUTH_TOKEN", "auth456")

	cfg, err := config.FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.TelegramToken != "tok123" {
		t.Errorf("TelegramToken = %q, want %q", cfg.TelegramToken, "tok123")
	}
	if cfg.ThingsAuthToken != "auth456" {
		t.Errorf("ThingsAuthToken = %q, want %q", cfg.ThingsAuthToken, "auth456")
	}
}

func TestFromEnv_MissingTelegramToken(t *testing.T) {
	t.Setenv("TELETHINGS_TELEGRAM_TOKEN", "")
	t.Setenv("TELETHINGS_THINGS_AUTH_TOKEN", "auth456")

	_, err := config.FromEnv()
	if !errors.Is(err, config.ErrMissingTelegramToken) {
		t.Errorf("expected ErrMissingTelegramToken, got %v", err)
	}
}

func TestFromEnv_MissingThingsAuthToken(t *testing.T) {
	t.Setenv("TELETHINGS_TELEGRAM_TOKEN", "tok123")
	t.Setenv("TELETHINGS_THINGS_AUTH_TOKEN", "")

	_, err := config.FromEnv()
	if !errors.Is(err, config.ErrMissingThingsAuthToken) {
		t.Errorf("expected ErrMissingThingsAuthToken, got %v", err)
	}
}
