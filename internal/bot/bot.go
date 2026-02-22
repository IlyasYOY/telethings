package bot

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/IlyasYOY/telethings/internal/config"
)

// Bot is the Telegram bot.
type Bot struct {
	api     *tgbotapi.BotAPI
	handler *Handler
}

// apiSender wraps tgbotapi.BotAPI to implement MessageSender.
type apiSender struct {
	api *tgbotapi.BotAPI
}

func (s *apiSender) Send(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := s.api.Send(msg)
	return err
}

// New creates a Bot from configuration and an Opener.
func New(cfg *config.Config, o opener) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, fmt.Errorf("create bot API: %w", err)
	}

	sender := &apiSender{api: api}
	handler := NewHandler(sender, o, cfg.ThingsAuthToken)

	return &Bot{api: api, handler: handler}, nil
}

// Run starts long-polling and processes updates until ctx is cancelled.
func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	log.Printf("bot started as @%s", b.api.Self.UserName)

	for {
		select {
		case <-ctx.Done():
			b.api.StopReceivingUpdates()
			return ctx.Err()
		case update, ok := <-updates:
			if !ok {
				return nil
			}
			if err := b.handler.Handle(update); err != nil {
				log.Printf("handle update: %v", err)
			}
		}
	}
}
