package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type webhookEmbed struct {
	Title       string `json:"title, omitempty"`
	Description string `json:"description, omitempty"`
	Color       int    `json:"color, omitempty"`
}

type webhookPayload struct {
	Embeds []webhookEmbed `json:"embeds"`
}

func SendReminder(webhookURL string, events []CalendarEvent, targetDate time.Time) error {
	embed := webhookEmbed{
		Title: fmt.Sprintf("📅 明日 (%s) の予定", targetDate.Format("2006-01-02 (Mon)")),
		Color: 0x4285F4,
	}
	if len(events) == 0 {
		embed.Description = "明日の予定はありません。"
	} else {
		var sb strings.Builder
		for _, e := range events {
			if e.AllDay {
				sb.WriteString(fmt.Sprintf("• **終日** %s\n", e.Summary))
			} else {
				sb.WriteString(fmt.Sprintf("• **%s** %s\n", e.StartTime.Format("15:04"), e.Summary))
			}
			if e.Location != "" {
				sb.WriteString(fmt.Sprintf("　　📍 %s\n", e.Location))
			}
		}
		embed.Description = sb.String()
	}

	payload := webhookPayload{Embeds: []webhookEmbed{embed}}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ペイロードの生成に失敗しました: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("リクエストの作成に失敗しました: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Webhookの送信に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Webhookがエラーを返しました: status=%d", resp.StatusCode)
	}
	return nil
}
