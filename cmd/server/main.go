package main

import (
	"new-client-notification-bot/config"
	"new-client-notification-bot/internal/handlers"
	"new-client-notification-bot/internal/services"
	"new-client-notification-bot/pkg/logger"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	config.Init()

	cfg := config.NewBotConfig()
	logCfg := config.NewLogConfig()

	customLogger := logger.NewLogger(logCfg)

	telegramBotService, err := services.NewTelegramBotService(cfg, customLogger)
	if err != nil {
		customLogger.Fatal().Err(err).Msg("failed to create telegram bot service")
	}

	app := fiber.New()
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: customLogger,
	}))
	app.Use(recover.New())

	handlers.NewNotificationHandler(app, telegramBotService, customLogger)

	err = app.Listen(":3000")
	if err != nil {
		customLogger.Fatal().Err(err).Msg("failed to listen")
	}
}
