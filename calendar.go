package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/apognu/gocal"
)

type CalendarEvent struct {
	Summary   string
	StartTime time.Time
	AllDay    bool
	Location  string
}

func GetTomorrowEvents(icalURL string) ([]CalendarEvent, error) {
	resp, err := http.Get(icalURL)
	if err != nil {
		return nil, fmt.Errorf("ICALファイルの取得に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ICALファイルの取得に失敗しました: status=%d", resp.StatusCode)
	}

	now := time.Now()
	loc := now.Location()
	tomorrow := now.AddDate(0, 0, 1)
	start := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, loc)
	end := start.AddDate(0, 0, 1)

	parser := gocal.NewParser(resp.Body)
	parser.Start, parser.End = &start, &end
	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("icalの解析に失敗しました: %w", err)
	}

	var result []CalendarEvent
	for _, e := range parser.Events {
		summary := e.Summary
		if summary == "" {
			summary = "(タイトルなし)"
		}
		allDay := e.End.Sub(*e.Start) >= 24*time.Hour
		result = append(result, CalendarEvent{
			Summary:   summary,
			StartTime: e.Start.In(loc),
			AllDay:    allDay,
			Location:  e.Location,
		})
	}
	return result, nil
}
