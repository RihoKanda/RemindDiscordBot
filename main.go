package main

import (
	"log"
	"time"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	events, err := GetTomorrowEvents(cfg.ICALURL)
	if err != nil {
		log.Fatalf("予定の取得に失敗しました: %v", err)
	}

	targetDate := time.Now().AddDate(0, 0, 1)
	if err := SendReminder(cfg.DiscordWebhookURL, events, targetDate); err != nil {
		log.Fatalf("Discordへの送信に失敗しました: %v", err)
	}
	log.Printf("リマインドを送信しました(%d件の予定)\n", len(events))
}
