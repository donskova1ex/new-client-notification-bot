package logger

import (
	"bytes"
	"new-client-notification-bot/config"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestNewLogger_JSONFormat(t *testing.T) {
	cfg := &config.LogConfig{
		Level:  int(zerolog.InfoLevel),
		Format: "json",
	}

	logger := NewLogger(cfg)

	if logger == nil {
		t.Fatal("Expected logger, got nil")
	}

	if zerolog.GlobalLevel() != zerolog.InfoLevel {
		t.Errorf("Expected global level %v, got %v", zerolog.InfoLevel, zerolog.GlobalLevel())
	}

	var buf bytes.Buffer

	testLogger := zerolog.New(&buf).With().Timestamp().Logger()
	testLogger.Info().Msg("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected log output to contain 'test message', got: %s", output)
	}

	if !strings.Contains(output, "{") || !strings.Contains(output, "}") {
		t.Errorf("Expected JSON format, got: %s", output)
	}
}

func TestNewLogger_ConsoleFormat(t *testing.T) {
	cfg := &config.LogConfig{
		Level:  int(zerolog.DebugLevel),
		Format: "console",
	}

	logger := NewLogger(cfg)

	if logger == nil {
		t.Fatal("Expected logger, got nil")
	}

	if zerolog.GlobalLevel() != zerolog.DebugLevel {
		t.Errorf("Expected global level %v, got %v", zerolog.DebugLevel, zerolog.GlobalLevel())
	}
}

func TestNewLogger_DefaultFormat(t *testing.T) {
	cfg := &config.LogConfig{
		Level:  int(zerolog.WarnLevel),
		Format: "unknown_format", // Неизвестный формат должен использовать default
	}

	logger := NewLogger(cfg)

	if logger == nil {
		t.Fatal("Expected logger, got nil")
	}

	if zerolog.GlobalLevel() != zerolog.WarnLevel {
		t.Errorf("Expected global level %v, got %v", zerolog.WarnLevel, zerolog.GlobalLevel())
	}
}

func TestNewLogger_DifferentLevels(t *testing.T) {
	levels := []struct {
		name  string
		level int
	}{
		{"Trace", int(zerolog.TraceLevel)},
		{"Debug", int(zerolog.DebugLevel)},
		{"Info", int(zerolog.InfoLevel)},
		{"Warn", int(zerolog.WarnLevel)},
		{"Error", int(zerolog.ErrorLevel)},
		{"Fatal", int(zerolog.FatalLevel)},
		{"Panic", int(zerolog.PanicLevel)},
		{"Disabled", int(zerolog.Disabled)},
	}

	for _, tt := range levels {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.LogConfig{
				Level:  tt.level,
				Format: "json",
			}

			logger := NewLogger(cfg)

			if logger == nil {
				t.Fatal("Expected logger, got nil")
			}

			expectedLevel := zerolog.Level(tt.level)
			if zerolog.GlobalLevel() != expectedLevel {
				t.Errorf("Expected global level %v, got %v", expectedLevel, zerolog.GlobalLevel())
			}
		})
	}
}

func TestNewLogger_LogOutput(t *testing.T) {
	tests := []struct {
		name   string
		format string
		level  int
	}{
		{
			name:   "JSON format with Info level",
			format: "json",
			level:  int(zerolog.InfoLevel),
		},
		{
			name:   "Console format with Debug level",
			format: "console",
			level:  int(zerolog.DebugLevel),
		},
		{
			name:   "Default format with Warn level",
			format: "unknown",
			level:  int(zerolog.WarnLevel),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.LogConfig{
				Level:  tt.level,
				Format: tt.format,
			}

			logger := NewLogger(cfg)

			if logger == nil {
				t.Fatal("Expected logger, got nil")
			}

			var buf bytes.Buffer
			var testLogger zerolog.Logger

			if tt.format == "json" {
				testLogger = zerolog.New(&buf).With().Timestamp().Logger()
			} else {
				consoleWriter := zerolog.ConsoleWriter{Out: &buf}
				testLogger = zerolog.New(consoleWriter).With().Timestamp().Logger()
			}

			testLogger = testLogger.Level(zerolog.Level(tt.level))

			if tt.level <= int(zerolog.InfoLevel) {
				testLogger.Info().Str("test", "value").Msg("test message")
			} else {
				testLogger.Warn().Str("test", "value").Msg("test message")
			}

			output := buf.String()
			if output == "" {
				t.Error("Expected log output, got empty string")
			}

			if tt.format == "json" {
				if !strings.Contains(output, "test message") {
					t.Errorf("Expected log to contain 'test message', got: %s", output)
				}
				if !strings.Contains(output, "test") {
					t.Errorf("Expected log to contain 'test', got: %s", output)
				}
			} else {
				if !strings.Contains(output, "test message") {
					t.Errorf("Expected log to contain 'test message', got: %s", output)
				}
			}
		})
	}
}

func TestNewLogger_NilConfig(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Expected panic with nil config: %v", r)
		}
	}()

	logger := NewLogger(nil)

	if logger != nil {
		t.Error("Expected nil logger or panic with nil config")
	}
}

func TestNewLogger_GlobalLevelIsolation(t *testing.T) {
	originalLevel := zerolog.GlobalLevel()
	defer func() {
		zerolog.SetGlobalLevel(originalLevel)
	}()

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	cfg := &config.LogConfig{
		Level:  int(zerolog.ErrorLevel),
		Format: "json",
	}

	logger := NewLogger(cfg)

	if zerolog.GlobalLevel() != zerolog.ErrorLevel {
		t.Errorf("Expected global level to be %v, got %v", zerolog.ErrorLevel, zerolog.GlobalLevel())
	}

	if logger == nil {
		t.Fatal("Expected logger, got nil")
	}
}

func TestNewLogger_ConsoleWriterConfiguration(t *testing.T) {
	cfg := &config.LogConfig{
		Level:  int(zerolog.InfoLevel),
		Format: "console",
	}

	logger := NewLogger(cfg)

	if logger == nil {
		t.Fatal("Expected logger, got nil")
	}

	var buf bytes.Buffer
	consoleWriter := zerolog.ConsoleWriter{Out: &buf}
	testLogger := zerolog.New(consoleWriter).With().Timestamp().Logger()

	testLogger.Info().Str("key", "value").Msg("console test")

	output := buf.String()
	if output == "" {
		t.Error("Expected console output, got empty string")
	}

	if !strings.Contains(output, "console test") {
		t.Errorf("Expected output to contain 'console test', got: %s", output)
	}
}

func BenchmarkNewLogger(b *testing.B) {
	cfg := &config.LogConfig{
		Level:  int(zerolog.InfoLevel),
		Format: "json",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewLogger(cfg)
	}
}

func BenchmarkNewLogger_Console(b *testing.B) {
	cfg := &config.LogConfig{
		Level:  int(zerolog.InfoLevel),
		Format: "console",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewLogger(cfg)
	}
}

func TestNewLogger_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.LogConfig
		expectPanic bool
	}{
		{
			name: "empty format",
			config: &config.LogConfig{
				Level:  int(zerolog.InfoLevel),
				Format: "",
			},
			expectPanic: false,
		},
		{
			name: "negative level",
			config: &config.LogConfig{
				Level:  -1,
				Format: "json",
			},
			expectPanic: false, // zerolog должен обработать это
		},
		{
			name: "very high level",
			config: &config.LogConfig{
				Level:  999,
				Format: "json",
			},
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic but didn't get one")
					}
				}()
			}

			logger := NewLogger(tt.config)

			if !tt.expectPanic && logger == nil {
				t.Error("Expected logger but got nil")
			}
		})
	}
}
