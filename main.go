package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/api/calendar/v3"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	ctx := context.Background()

	calSrv, err := NewCalendarService(ctx, cfg)
	if err != nil {
		log.Fatalf("Google Calenderサービスの初期化に失敗しました: %v", err)
	}

	log.Printf("Botを起動しました 毎日 %02d:%02d に翌日の予定をチャンネル %s へ通知します \n",
		cfg.RemindHour, cfg.RemindMinute)

	go runScheduler(calSrv, cfg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Botを終了します")
}

func runScheduler(calSrv *calendar.Service, cfg *Config) {
	for {
		next := nextRunTime(cfg.RemindHour, cfg.RemindMinute)
		wait := time.Until(next)
		log.Printf("次回実行予定: %s (あと %s)\n", next.Format("2006-01-02 15:04:05"), wait.Round(time.Second))

		timer := time.NewTimer(wait)
		<-timer.C

		if err := executeReminder(calSrv, cfg); err != nil {
			log.Printf("リマインド実行中にエラーが発生しました: %v\n", err)
		}
	}
}

func nextRunTime(hour, minute int) time.Time {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	if !next.After(now) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}

func executeReminder(calSrv *calendar.Service, cfg *Config) error {
	events, err := GetTomorrowEvents(calSrv, cfg.GoogleCalendarID)
	if err != nil {
		return err
	}
	targetDate := time.Now().AddDate(0, 0, 1)
	if err := SendReminder(cfg.DiscordWebhookURL, events, targetDate); err != nil {
		return err
	}
	log.Printf("リマインドを送信しました(%d件の予定)\n", len(events))
	return nil
}
