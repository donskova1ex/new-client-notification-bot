package config

import (
	"os"
	"testing"
)

func TestNewBotConfig(t *testing.T) {
	originalBotToken := os.Getenv("BOT_TOKEN")
	originalChatID := os.Getenv("CHAT_ID")

	defer func() {
		if originalBotToken != "" {
			os.Setenv("BOT_TOKEN", originalBotToken)
		} else {
			os.Unsetenv("BOT_TOKEN")
		}
		if originalChatID != "" {
			os.Setenv("CHAT_ID", originalChatID)
		} else {
			os.Unsetenv("CHAT_ID")
		}
	}()

	tests := []struct {
		name        string
		botToken    string
		chatID      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid configuration",
			botToken:    "123456789:ABCdefGHIjklMNOpqrsTUVwxyz",
			chatID:      "123456789",
			expectError: false,
		},
		{
			name:        "empty bot token",
			botToken:    "",
			chatID:      "123456789",
			expectError: true,
			errorMsg:    "bot token required",
		},
		{
			name:        "empty chat ID",
			botToken:    "123456789:ABCdefGHIjklMNOpqrsTUVwxyz",
			chatID:      "",
			expectError: true,
			errorMsg:    "bot chat id required",
		},
		{
			name:        "invalid chat ID",
			botToken:    "123456789:ABCdefGHIjklMNOpqrsTUVwxyz",
			chatID:      "invalid",
			expectError: true,
			errorMsg:    "bot chat id required",
		},
		{
			name:        "zero chat ID",
			botToken:    "123456789:ABCdefGHIjklMNOpqrsTUVwxyz",
			chatID:      "0",
			expectError: true,
			errorMsg:    "bot chat id required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.botToken != "" {
				os.Setenv("BOT_TOKEN", tt.botToken)
			} else {
				os.Unsetenv("BOT_TOKEN")
			}

			if tt.chatID != "" {
				os.Setenv("CHAT_ID", tt.chatID)
			} else {
				os.Unsetenv("CHAT_ID")
			}

			cfg, err := NewBotConfig()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q, got %q", tt.errorMsg, err.Error())
				}
				if cfg != nil {
					t.Errorf("expected nil config on error, got %+v", cfg)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if cfg == nil {
					t.Errorf("expected config, got nil")
				} else {
					if cfg.BotToken != tt.botToken {
						t.Errorf("expected bot token %q, got %q", tt.botToken, cfg.BotToken)
					}
					if cfg.ChatID != 123456789 { // Ожидаем int64 значение
						t.Errorf("expected chat ID %d, got %d", 123456789, cfg.ChatID)
					}
				}
			}
		})
	}
}

func TestNewLogConfig(t *testing.T) {
	originalLogLevel := os.Getenv("LOG_LEVEL")
	originalLogFormat := os.Getenv("LOG_FORMAT")

	defer func() {
		if originalLogLevel != "" {
			os.Setenv("LOG_LEVEL", originalLogLevel)
		} else {
			os.Unsetenv("LOG_LEVEL")
		}
		if originalLogFormat != "" {
			os.Setenv("LOG_FORMAT", originalLogFormat)
		} else {
			os.Unsetenv("LOG_FORMAT")
		}
	}()

	tests := []struct {
		name           string
		logLevel       string
		logFormat      string
		expectedLevel  int
		expectedFormat string
	}{
		{
			name:           "default values",
			logLevel:       "",
			logFormat:      "",
			expectedLevel:  0,
			expectedFormat: "json",
		},
		{
			name:           "custom values",
			logLevel:       "1",
			logFormat:      "text",
			expectedLevel:  1,
			expectedFormat: "text",
		},
		{
			name:           "invalid log level",
			logLevel:       "invalid",
			logFormat:      "json",
			expectedLevel:  0, // Должен вернуться к дефолтному значению
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.logLevel != "" {
				os.Setenv("LOG_LEVEL", tt.logLevel)
			} else {
				os.Unsetenv("LOG_LEVEL")
			}

			if tt.logFormat != "" {
				os.Setenv("LOG_FORMAT", tt.logFormat)
			} else {
				os.Unsetenv("LOG_FORMAT")
			}

			cfg := NewLogConfig()

			if cfg == nil {
				t.Errorf("expected config, got nil")
			} else {
				if cfg.Level != tt.expectedLevel {
					t.Errorf("expected log level %d, got %d", tt.expectedLevel, cfg.Level)
				}
				if cfg.Format != tt.expectedFormat {
					t.Errorf("expected log format %q, got %q", tt.expectedFormat, cfg.Format)
				}
			}
		})
	}
}

func TestGetString(t *testing.T) {
	originalValue := os.Getenv("TEST_STRING_VAR")
	defer func() {
		if originalValue != "" {
			os.Setenv("TEST_STRING_VAR", originalValue)
		} else {
			os.Unsetenv("TEST_STRING_VAR")
		}
	}()

	tests := []struct {
		name         string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "existing environment variable",
			envValue:     "test_value",
			defaultValue: "default",
			expected:     "test_value",
		},
		{
			name:         "empty environment variable",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "unset environment variable",
			envValue:     "", // Не устанавливаем переменную
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("TEST_STRING_VAR", tt.envValue)
			} else {
				os.Unsetenv("TEST_STRING_VAR")
			}

			result := getString("TEST_STRING_VAR", tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	originalValue := os.Getenv("TEST_INT_VAR")
	defer func() {
		if originalValue != "" {
			os.Setenv("TEST_INT_VAR", originalValue)
		} else {
			os.Unsetenv("TEST_INT_VAR")
		}
	}()

	tests := []struct {
		name         string
		envValue     string
		defaultValue int
		expected     int
	}{
		{
			name:         "valid integer",
			envValue:     "42",
			defaultValue: 0,
			expected:     42,
		},
		{
			name:         "invalid integer",
			envValue:     "not_a_number",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "empty value",
			envValue:     "",
			defaultValue: 5,
			expected:     5,
		},
		{
			name:         "zero value",
			envValue:     "0",
			defaultValue: 10,
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("TEST_INT_VAR", tt.envValue)
			} else {
				os.Unsetenv("TEST_INT_VAR")
			}

			result := getInt("TEST_INT_VAR", tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

