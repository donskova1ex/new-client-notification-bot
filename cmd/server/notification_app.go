package main

import (
	"errors"
	"net/http"
	"new-client-notification-bot/config"
	"new-client-notification-bot/internal/handlers"
	"new-client-notification-bot/internal/services"
	"new-client-notification-bot/pkg/logger"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	config.Init()

	logCfg := config.NewLogConfig()

	customLogger := logger.NewLogger(logCfg)

	cfg, err := config.NewBotConfig()
	if err != nil {
		customLogger.Fatal().Err(err).Msg("Failed to load bot config")
	}

	telegramBotService, err := services.NewTelegramBotService(cfg, customLogger)
	if err != nil {
		customLogger.Fatal().Err(err).Msg("failed to create telegram bot service")
	}

	app := fiber.New()
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: customLogger,
	}))
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "POST",
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "too many requests",
			})
		},
	}))

	handlers.NewNotificationHandler(app, telegramBotService, customLogger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		if err := app.Listen(":3000"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			customLogger.Fatal().Err(err).Msg("failed to listen")
		}
	}()

	<-c
	customLogger.Info().Msg("shutting down")
	if err := app.Shutdown(); err != nil {
		customLogger.Fatal().Err(err).Msg("failed to shutdown")
	}

}
