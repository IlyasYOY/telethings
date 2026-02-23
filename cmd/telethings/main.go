package main

import (
	"context"
	"errors"
	"log"
	"os/signal"
	"syscall"

	"github.com/IlyasYOY/telethings/internal/bot"
	"github.com/IlyasYOY/telethings/internal/config"
	"github.com/IlyasYOY/telethings/internal/opener"
	"github.com/IlyasYOY/telethings/internal/thingsreader"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	b, err := bot.New(cfg, opener.MacOSOpener{}, thingsreader.AppleScriptReader{})
	if err != nil {
		log.Fatalf("bot: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := b.Run(ctx); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		log.Printf("bot stopped: %v", err)
	}
}
