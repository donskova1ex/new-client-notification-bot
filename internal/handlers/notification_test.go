package handlers

import (
	"new-client-notification-bot/internal/domain"
	"testing"
)

func TestEmptyStringValidator(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		fieldName   string
		expectError bool
	}{
		{
			name:        "valid non-empty string",
			input:       "test value",
			fieldName:   "test_field",
			expectError: false,
		},
		{
			name:        "empty string should return error",
			input:       "",
			fieldName:   "test_field",
			expectError: true,
		},
		{
			name:        "whitespace only string should return error",
			input:       "   ",
			fieldName:   "test_field",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := emptyStringValidator(tt.input, tt.fieldName)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestPhoneValidation(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected bool
	}{
		{
			name:     "valid phone with +7",
			phone:    "+7 912 345 67 89",
			expected: true,
		},
		{
			name:     "valid phone with 8",
			phone:    "8 912 345 67 89",
			expected: true,
		},
		{
			name:     "valid phone with 7",
			phone:    "7 912 345 67 89",
			expected: true,
		},
		{
			name:     "valid phone without spaces",
			phone:    "79123456789",
			expected: true,
		},
		{
			name:     "invalid phone - wrong operator",
			phone:    "+7 812 345 67 89",
			expected: false,
		},
		{
			name:     "invalid phone - too short",
			phone:    "+7 912 345 67",
			expected: false,
		},
		{
			name:     "invalid phone - too long",
			phone:    "+7 912 345 67 89 12",
			expected: false,
		},
		{
			name:     "invalid phone - empty",
			phone:    "",
			expected: false,
		},
		{
			name:     "invalid phone - letters",
			phone:    "+7 912 abc 67 89",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := phoneValidation(tt.phone)
			if result != tt.expected {
				t.Errorf("phoneValidation(%q) = %v, expected %v", tt.phone, result, tt.expected)
			}
		})
	}
}

func TestCreateFormatNotification(t *testing.T) {
	handler := &Notification{}

	tests := []struct {
		name     string
		input    domain.Notification
		expected string
	}{
		{
			name: "valid notification formatting",
			input: domain.Notification{
				Phone:            "+7 912 345 67 89",
				CompanyName:      "Test Company",
				NotificationText: "Test message",
			},
			expected: "Клиент: Test Company;\nТелефон: +7 912 345 67 89;\nТекст обращение: Test message",
		},
		{
			name: "notification with special characters",
			input: domain.Notification{
				Phone:            "8 912 345 67 89",
				CompanyName:      "ООО \"Рога и копыта\"",
				NotificationText: "Сообщение с переносами\nстрок",
			},
			expected: "Клиент: ООО \"Рога и копыта\";\nТелефон: 8 912 345 67 89;\nТекст обращение: Сообщение с переносами\nстрок",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.createFormatNotification(&tt.input)
			if result != tt.expected {
				t.Errorf("createFormatNotification() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestValidateRequest(t *testing.T) {
	handler := &Notification{}

	tests := []struct {
		name        string
		input       domain.Notification
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid request",
			input: domain.Notification{
				Phone:            "+7 912 345 67 89",
				CompanyName:      "Test Company",
				NotificationText: "Test message",
			},
			expectError: false,
		},
		{
			name: "empty phone",
			input: domain.Notification{
				Phone:            "",
				CompanyName:      "Test Company",
				NotificationText: "Test message",
			},
			expectError: true,
			errorMsg:    "phone is required",
		},
		{
			name: "empty company name",
			input: domain.Notification{
				Phone:            "+7 912 345 67 89",
				CompanyName:      "",
				NotificationText: "Test message",
			},
			expectError: true,
			errorMsg:    "company_name is required",
		},
		{
			name: "empty notification text",
			input: domain.Notification{
				Phone:            "+7 912 345 67 89",
				CompanyName:      "Test Company",
				NotificationText: "",
			},
			expectError: true,
			errorMsg:    "notification_text is required",
		},
		{
			name: "invalid phone",
			input: domain.Notification{
				Phone:            "invalid phone",
				CompanyName:      "Test Company",
				NotificationText: "Test message",
			},
			expectError: true,
			errorMsg:    "invalid phone",
		},
		{
			name: "notification text too long",
			input: domain.Notification{
				Phone:            "+7 912 345 67 89",
				CompanyName:      "Test Company",
				NotificationText: string(make([]byte, 256)), // 256 символов
			},
			expectError: true,
			errorMsg:    "notification_text too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.validateRequest(&tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

