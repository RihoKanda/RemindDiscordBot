package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Botの動作に必要な設定のまとめ
type Config struct {
	ICALURL           string // Google CalendarのICAL URL
	DiscordWebhookURL string // Discord WebhookのURL
}

// .envと環境変数から設定の読み込み
func LoadConfig() (*Config, error) {
	// .envファイルがあれば読み込む
	if err := godotenv.Load(); err != nil {
		log.Println(".env ファイルがないよ　環境変数から読み込むよ")
	}

	cfg := &Config{
		ICALURL:           os.Getenv("ICAL_URL"),
		DiscordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
	}

	if cfg.ICALURL == "" {
		return nil, fmt.Errorf("ICAL_URL が設定されていないよ")
	}
	if cfg.DiscordWebhookURL == "" {
		return nil, fmt.Errorf("DISCORD_WEBHOOK_URL が設定されていないよ")
	}

	return cfg, nil
}
