package logger

import (
	"new-client-notification-bot/config"
	"os"

	"github.com/rs/zerolog"
)

func NewLogger(cfg *config.LogConfig) *zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.Level(cfg.Level))

	var logger zerolog.Logger
	switch cfg.Format {
	case "json":
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	default:
		consoleWritter := zerolog.ConsoleWriter{Out: os.Stdout}
		logger = zerolog.New(consoleWritter).With().Timestamp().Logger()
	}
	return &logger

}
