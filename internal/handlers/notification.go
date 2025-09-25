package handlers

import (
	"fmt"
	"new-client-notification-bot/internal/domain"
	"new-client-notification-bot/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Notification struct {
	router             fiber.Router
	telegramBotService *services.TelegramBotService
	logger             *zerolog.Logger
}

func NewNotificationHandler(router fiber.Router, telegramBotService *services.TelegramBotService, logger *zerolog.Logger) {
	handler := &Notification{
		router:             router,
		telegramBotService: telegramBotService,
		logger:             logger,
	}
	api := handler.router.Group("/api/v1")
	api.Post("/notification", handler.CreateNotification)

}

func (n *Notification) CreateNotification(c *fiber.Ctx) error {
	var req domain.Notification
	n.logger.Info().Str("ip", c.IP()).Msg("received request")
	if err := c.BodyParser(&req); err != nil {
		n.logger.Error().Err(err).Msg("failed to parse request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "failed to parse request",
		})
	}

	if err := n.validateRequest(&req); err != nil {
		n.logger.Error().Err(err).Msg("failed to validate request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "failed to validate request",
		})
	}

	message := n.createFormatNotification(&req)

	if err := n.telegramBotService.SendMessage(c.Context(), message); err != nil {
		n.logger.Error().Err(err).Msg("failed to send message")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "failed to send message",
		})
	}

	n.logger.Info().Interface("request", req).Msg("sent message successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "sent message successfully",
	})
}

func (n *Notification) validateRequest(req *domain.Notification) error {
	if err := emptyStringValidator(req.Phone, "phone"); err != nil {
		return err
	}

	if err := emptyStringValidator(req.CompanyName, "company_name"); err != nil {
		return err
	}

	if err := emptyStringValidator(req.NotificationText, "notification_text"); err != nil {
		return err
	}

	return nil
}

func emptyStringValidator(s, stringName string) error {
	if s == "" {
		return fmt.Errorf("%s is required", stringName)
	}
	return nil
}

func (n *Notification) createFormatNotification(req *domain.Notification) string {
	formatMessage := fmt.Sprintf(
		"Клиент: %s;\nТелефон: %s;\nТекст обращение: %s",
		req.Phone,
		req.CompanyName,
		req.NotificationText,
	)
	return formatMessage
}
