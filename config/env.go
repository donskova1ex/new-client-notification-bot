package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type BotConfig struct {
	BotToken string
	ChatID   int64
}

type LogConfig struct {
	Level  int
	Format string
}

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("successfully loaded .env file")
}

func getInt(key string, defaultString int) int {
	valStr := os.Getenv(key)

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultString
	}

	return val
}

func getString(key, defaultString string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultString
	}
	return val
}

func NewBotConfig() *BotConfig {
	return &BotConfig{
		BotToken: getString("BOT_TOKEN", ""),
		ChatID:   int64(getInt("CHAT_ID", 0)),
	}
}

func NewLogConfig() *LogConfig {
	return &LogConfig{
		Level:  getInt("LOG_LEVEL", 0),
		Format: getString("LOG_FORMAT", "json"),
	}
}
