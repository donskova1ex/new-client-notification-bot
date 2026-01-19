package services

import "context"

type TelegramBotServiceInterface interface {
	SendMessage(ctx context.Context, message string) error
}

