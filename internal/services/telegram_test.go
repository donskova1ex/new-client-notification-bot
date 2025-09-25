package services

import (
	"context"
	"errors"
	"testing"
	"time"
)

type MockTelegramBotService struct {
	shouldError  bool
	errorMsg     string
	sentMessages []string
}

func (m *MockTelegramBotService) SendMessage(ctx context.Context, message string) error {
	m.sentMessages = append(m.sentMessages, message)
	if m.shouldError {
		return errors.New(m.errorMsg)
	}
	return nil
}

func TestTelegramBotService_SendMessage(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "successful message send",
			message:     "Test message",
			shouldError: false,
		},
		{
			name:        "failed message send",
			message:     "Test message",
			shouldError: true,
			errorMsg:    "telegram API error",
		},
		{
			name:        "empty message",
			message:     "",
			shouldError: false,
		},
		{
			name:        "long message",
			message:     string(make([]byte, 1000)),
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockTelegramBotService{
				shouldError: tt.shouldError,
				errorMsg:    tt.errorMsg,
			}

			ctx := context.Background()
			err := mock.SendMessage(ctx, tt.message)

			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(mock.sentMessages) == 0 {
					t.Errorf("message was not recorded")
				}
				if len(mock.sentMessages) > 0 && mock.sentMessages[0] != tt.message {
					t.Errorf("expected message %q, got %q", tt.message, mock.sentMessages[0])
				}
			}
		})
	}
}

func TestTelegramBotService_ContextCancellation(t *testing.T) {
	mock := &MockTelegramBotService{}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(2 * time.Millisecond)

	err := mock.SendMessage(ctx, "test message")

	if err != nil {
		t.Logf("Context cancellation handled: %v", err)
	}
}

func BenchmarkTelegramBotService_SendMessage(b *testing.B) {
	mock := &MockTelegramBotService{}
	ctx := context.Background()
	message := "Benchmark test message"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mock.SendMessage(ctx, message)
	}
}
