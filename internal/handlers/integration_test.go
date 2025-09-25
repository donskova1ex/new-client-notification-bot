package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"new-client-notification-bot/internal/domain"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type MockTelegramService struct {
	shouldError  bool
	errorMsg     string
	sentMessages []string
}

func (m *MockTelegramService) SendMessage(ctx context.Context, message string) error {
	m.sentMessages = append(m.sentMessages, message)
	if m.shouldError {
		return errors.New(m.errorMsg)
	}
	return nil
}

func setupTestApp(telegramService *MockTelegramService) *fiber.App {
	app := fiber.New()

	logger := zerolog.Nop()

	handler := &Notification{
		router:             app,
		telegramBotService: telegramService,
		logger:             &logger,
	}

	api := app.Group("/api/v1")
	api.Post("/notification", handler.CreateNotification)

	return app
}

func TestCreateNotification_Integration(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      domain.Notification
		telegramError    bool
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{
		{
			name: "successful notification",
			requestBody: domain.Notification{
				Phone:            "+7 912 345 67 89",
				CompanyName:      "Test Company",
				NotificationText: "Test message",
			},
			telegramError:  false,
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"success": true,
				"message": "sent message successfully",
			},
		},
		{
			name: "invalid phone number",
			requestBody: domain.Notification{
				Phone:            "invalid phone",
				CompanyName:      "Test Company",
				NotificationText: "Test message",
			},
			telegramError:  false,
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"success": false,
				"message": "failed to validate request",
			},
		},
		{
			name: "empty phone",
			requestBody: domain.Notification{
				Phone:            "",
				CompanyName:      "Test Company",
				NotificationText: "Test message",
			},
			telegramError:  false,
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"success": false,
				"message": "failed to validate request",
			},
		},
		{
			name: "telegram service error",
			requestBody: domain.Notification{
				Phone:            "+7 912 345 67 89",
				CompanyName:      "Test Company",
				NotificationText: "Test message",
			},
			telegramError:  true,
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"success": false,
				"message": "failed to send message",
			},
		},
		{
			name: "notification text too long",
			requestBody: domain.Notification{
				Phone:            "+7 912 345 67 89",
				CompanyName:      "Test Company",
				NotificationText: string(make([]byte, 256)), // 256 символов
			},
			telegramError:  false,
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"success": false,
				"message": "failed to validate request",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTelegram := &MockTelegramService{
				shouldError: tt.telegramError,
				errorMsg:    "telegram API error",
			}

			app := setupTestApp(mockTelegram)

			jsonBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/notification", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to execute request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			var response map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if response["success"] != tt.expectedResponse["success"] {
				t.Errorf("Expected success %v, got %v", tt.expectedResponse["success"], response["success"])
			}

			if response["message"] != tt.expectedResponse["message"] {
				t.Errorf("Expected message %q, got %q", tt.expectedResponse["message"], response["message"])
			}

			if tt.expectedStatus == http.StatusOK {
				if len(mockTelegram.sentMessages) == 0 {
					t.Errorf("Expected message to be sent, but none was sent")
				} else {
					expectedMessage := "Клиент: Test Company;\nТелефон: +7 912 345 67 89;\nТекст обращение: Test message"
					if mockTelegram.sentMessages[0] != expectedMessage {
						t.Errorf("Expected message %q, got %q", expectedMessage, mockTelegram.sentMessages[0])
					}
				}
			}
		})
	}
}

func TestCreateNotification_InvalidJSON(t *testing.T) {
	mockTelegram := &MockTelegramService{}
	app := setupTestApp(mockTelegram)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/notification", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["success"] != false {
		t.Errorf("Expected success false, got %v", response["success"])
	}

	if response["message"] != "failed to parse request" {
		t.Errorf("Expected message 'failed to parse request', got %q", response["message"])
	}
}

func TestCreateNotification_WrongMethod(t *testing.T) {
	mockTelegram := &MockTelegramService{}
	app := setupTestApp(mockTelegram)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/notification", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}
