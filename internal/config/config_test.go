package config_test

import (
	"errors"
	"testing"

	"github.com/IlyasYOY/telethings/internal/config"
)

func TestFromEnv_Success(t *testing.T) {
	t.Setenv("TELETHINGS_TELEGRAM_TOKEN", "tok123")
	t.Setenv("TELETHINGS_THINGS_AUTH_TOKEN", "auth456")
	t.Setenv("TELETHINGS_ALLOWED_USER_IDS", "123,456,789")

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
	if len(cfg.AllowedUserIDs) != 3 {
		t.Errorf("AllowedUserIDs count = %d, want 3", len(cfg.AllowedUserIDs))
	}
	expectedIDs := []int64{123, 456, 789}
	for i, id := range cfg.AllowedUserIDs {
		if id != expectedIDs[i] {
			t.Errorf("AllowedUserIDs[%d] = %d, want %d", i, id, expectedIDs[i])
		}
	}
}

func TestFromEnv_MissingTelegramToken(t *testing.T) {
	t.Setenv("TELETHINGS_TELEGRAM_TOKEN", "")
	t.Setenv("TELETHINGS_THINGS_AUTH_TOKEN", "auth456")
	t.Setenv("TELETHINGS_ALLOWED_USER_IDS", "123")

	_, err := config.FromEnv()
	if !errors.Is(err, config.ErrMissingTelegramToken) {
		t.Errorf("expected ErrMissingTelegramToken, got %v", err)
	}
}

func TestFromEnv_MissingThingsAuthToken(t *testing.T) {
	t.Setenv("TELETHINGS_TELEGRAM_TOKEN", "tok123")
	t.Setenv("TELETHINGS_THINGS_AUTH_TOKEN", "")
	t.Setenv("TELETHINGS_ALLOWED_USER_IDS", "123")

	_, err := config.FromEnv()
	if !errors.Is(err, config.ErrMissingThingsAuthToken) {
		t.Errorf("expected ErrMissingThingsAuthToken, got %v", err)
	}
}

func TestFromEnv_MissingAllowedUserIDs(t *testing.T) {
	t.Setenv("TELETHINGS_TELEGRAM_TOKEN", "tok123")
	t.Setenv("TELETHINGS_THINGS_AUTH_TOKEN", "auth456")
	t.Setenv("TELETHINGS_ALLOWED_USER_IDS", "")

	_, err := config.FromEnv()
	if !errors.Is(err, config.ErrMissingAllowedUserIDs) {
		t.Errorf("expected ErrMissingAllowedUserIDs, got %v", err)
	}
}

func TestFromEnv_InvalidAllowedUserIDs(t *testing.T) {
	t.Setenv("TELETHINGS_TELEGRAM_TOKEN", "tok123")
	t.Setenv("TELETHINGS_THINGS_AUTH_TOKEN", "auth456")
	t.Setenv("TELETHINGS_ALLOWED_USER_IDS", "123,invalid,456")

	_, err := config.FromEnv()
	if !errors.Is(err, config.ErrInvalidAllowedUserIDs) {
		t.Errorf("expected ErrInvalidAllowedUserIDs, got %v", err)
	}
}

func TestFromEnv_AllowedUserIDsWithWhitespace(t *testing.T) {
	t.Setenv("TELETHINGS_TELEGRAM_TOKEN", "tok123")
	t.Setenv("TELETHINGS_THINGS_AUTH_TOKEN", "auth456")
	t.Setenv("TELETHINGS_ALLOWED_USER_IDS", " 123 , 456 , 789 ")

	cfg, err := config.FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.AllowedUserIDs) != 3 {
		t.Errorf("AllowedUserIDs count = %d, want 3", len(cfg.AllowedUserIDs))
	}
}
