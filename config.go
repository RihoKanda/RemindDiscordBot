package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Botの動作に必要な設定のまとめ
type Config struct {
	DiscordBotToken  string // DiscordBotのとーくん
	DiscordChannelID string // リマインドを投稿するチャンネルID
	GoogleCalendarID string // 対象のGoogleカレンダーID
	RemindHour       int    // リマインドを実行する時刻　0-23
	RemindMinute     int    // リマインドを実行する時刻　0-59
	CredentialsPath  string // GoogleOAuth2 jsonのパス
	TokenPath        string // 取得したOAuth2トークンの保存先パス
}

// .envと環境変数から設定の読み込み
func LoadConfig() (*Config, error) {
	// .envファイルがあれば読み込む
	if err := godotenv.Load(); err != nil {
		log.Println(".env ファイルがないよ　環境変数から読み込むよ")
	}

	cfg := &Config{
		DiscordBotToken:  os.Getenv("DISCORD_BOT_TOKEN"),
		DiscordChannelID: os.Getenv("DISCORD_CHANNEL_ID"),
		GoogleCalendarID: getEnvOrDefault("GOOGLE_CALENDAR_ID", "primary"),
		CredentialsPath:  getEnvOrDefault("GOOGLE_CREDENTIALS_PATH", "credentials.json"),
		TokenPath:        getEnvOrDefault("GOOGLE_TOKEN_PATH", "token.json"),
	}

	hourStr := getEnvOrDefault("REMIND_HOUR", "20")
	minuteStr := getEnvOrDefault("REMIND_MINUTE", "0")

	hour, err := strconv.Atoi(hourStr)
	if err != nil || hour < 0 || hour > 23 {
		return nil, fmt.Errorf("REMIND_HOUR が不正だよ: %s", hourStr)
	}
	minute, err := strconv.Atoi(minuteStr)
	if err != nil || minute < 0 || minute > 59 {
		return nil, fmt.Errorf("REMIND_MINUTE が不正だよ: %s", minuteStr)
	}
	cfg.RemindHour = hour
	cfg.RemindMinute = minute

	if cfg.DiscordBotToken == "" {
		return nil, fmt.Errorf("DISCORD_BOT_TOKEN が設定されていないよ")
	}
	if cfg.DiscordChannelID == "" {
		return nil, fmt.Errorf("DISCORD_CHANNEL_ID が設定されていないよ")
	}

	return cfg, nil
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}