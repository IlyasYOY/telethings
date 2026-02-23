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

func (s *apiSender) SendWithInlineKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	_, err := s.api.Send(msg)
	return err
}

func (s *apiSender) SendTyping(chatID int64) error {
	_, err := s.api.Request(tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping))
	return err
}

func (s *apiSender) AckCallback(callbackID string) error {
	_, err := s.api.Request(tgbotapi.NewCallback(callbackID, ""))
	return err
}

// New creates a Bot from configuration and an Opener.
func New(cfg *config.Config, o opener, r thingsReader) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, fmt.Errorf("create bot API: %w", err)
	}

	// Register commands with Telegram
	commands := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{Command: "start", Description: "Welcome message and quick help"},
		tgbotapi.BotCommand{Command: "add", Description: "Add a task to Things 3"},
		tgbotapi.BotCommand{Command: "today", Description: "Show today's tasks from Things 3"},
		tgbotapi.BotCommand{Command: "inbox", Description: "Show your Things 3 inbox"},
		tgbotapi.BotCommand{Command: "anytime", Description: "Show Anytime tasks (paged)"},
		tgbotapi.BotCommand{Command: "someday", Description: "Show Someday tasks (paged)"},
		tgbotapi.BotCommand{Command: "tags", Description: "Show tags and browse tasks by tag"},
	)
	if _, err := api.Request(commands); err != nil {
		return nil, fmt.Errorf("set bot commands: %w", err)
	}

	sender := &apiSender{api: api}
	handler := NewHandler(sender, o, r, cfg.ThingsAuthToken, cfg.AllowedUserIDs)

	return &Bot{api: api, handler: handler}, nil
}

// Run starts long-polling and processes updates until ctx is cancelled.
func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	// Keep long polling enabled to reduce request churn while waiting for updates.
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
