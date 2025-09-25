package services

import (
	"context"
	"net/http"
	"new-client-notification-bot/config"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type TelegramBotService struct {
	bot    *tgbotapi.BotAPI
	chatID int64
	logger *zerolog.Logger
}

func NewTelegramBotService(cfg *config.BotConfig, logger *zerolog.Logger) (*TelegramBotService, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create telegram bot")
		return nil, err
	}
	bot.Client = &http.Client{
		Timeout: 30 * time.Second,
	}
	logger.Info().Str("bot_name", bot.Self.UserName).Msg("telegram bot created")

	return &TelegramBotService{
		bot:    bot,
		chatID: cfg.ChatID,
		logger: logger,
	}, nil
}

func (t *TelegramBotService) SendMessage(ctx context.Context, message string) error {
	t.logger.Info().Int64("chat_id", t.chatID).Msg("sending message")

	msg := tgbotapi.NewMessage(t.chatID, message)

	if _, err := t.bot.Send(msg); err != nil {
		t.logger.Error().Err(err).Msg("failed to send message")
		return err
	}

	t.logger.Info().Int64("chat_id", t.chatID).Msg("message sent")
	return nil
}
