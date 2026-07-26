package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func SendReminder(session *discordgo.Session, channelID string, events []CalenderEvent, targetDate  time.Time) error {
	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("📅 明日 (%s) の予定", targetDate.Format("2006/01/02 (Mon)")),
		Color: 0x4285F4,
	}

	if len(events) == 0 {
		embed.Description = "明日の予定はありません"
	} else {
		var sb strings.Builder
		for _, e := renge events {
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

	_, err := session.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		return fmt.Errorf("Discordへの送信に失敗しました: %w", err)
	}
	return nil
}