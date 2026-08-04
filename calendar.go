package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func NewCalenderService(ctx context.Context, cfg *Config) (*calendar.Service, error) {
	b, err := os.ReadFile(cfg.CredentialsPath)
	if err != nil {
		return nil, fmt.Errorf("credentials.json の読み込みに失敗しました (%s): %w", cfg.CredentialsPath, err)
	}

	oauthConfig, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("credentials.json のパースに失敗しました: %w", err)
	}

	client, err := getHTTPClient(oauthConfig, cfg.TokenPath)
	if err != nil {
		return nil, err
	}

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("Calendarサービスの作成に失敗しました: %w", err)
	}
	return srv, nil
}

func getHTTPClient(config *oauth2.Config, tokenPath string) (*http.Client, error) {
	tok, err := tokenFromFile(tokenPath)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		if err := saveToken(tokenPath, tok); err != nil {
			return nil, err
		}
	}

	ts := config.TokenSource(context.Background(), tok)
	return oauth2.NewClient(context.Background(), tokenSaver{ts, tokenPath}), nil
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Println("=== Google Calendar 初回認証 ===")
	fmt.Println("以下のURLをブラウザで開き、Googleアカウントで許可してください:")
	fmt.Println(authURL)
	fmt.Print("表示された認可コードを貼り付けてEnterを押してください: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("認可コードの読み取りに失敗しました: %w", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, fmt.Errorf("トークンの取得に失敗しました: %w", err)
	}
	return tok, nil
}

func tokenFromFile(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	if err := json.NewDecoder(f).Decode(tok); err != nil {
		return nil, err
	}
	return tok, nil
}

func saveToken(path string, token *oauth2.Token) error {
	log.Printf("トークンを保存します: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("トークンファイルの保存に失敗しました: %w", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

type tokenSaver struct {
	src       oauth2.TokenSource
	tokenPath string
}

func (t tokenSaver) RoundTrip(req *http.Request) (*http.Response, error) {
	tok, err := t.src.Token()
	if err != nil {
		return nil, err
	}

	_ = saveToken(t.tokenPath, tok)
	req2 := req.Clone(req.Context())
	tok.SetAuthHeader(req2)
	return http.DefaultTransport.RoundTrip(req2)
}

type CalenderEvent struct {
	Summary   string
	StartTime time.Time
	AllDay    bool
	Location  string
}

func GetTomorrowEvents(srv *calendar.Service, calendarID string) ([]CalendarEvent, error) {
	now := time.Now()
	loc := now.Location()

	tomorrow := now.AddDate(0, 0, 1)
	startOfDay := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	events, err := srv.Events.List(calendarID).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(startOfDay.Format(time.RFC3339)).
		TimeMax(endOfDay.Format(time.RFC3339)).
		OrderBy("startTime").
		Do()
	if err != nil {
		return nil, fmt.Errorf("予定の取得に失敗しました: %w", err)
	}

	var result []CalendarEvent
	for _, item := range events.Items {
		ce := CalendarEvent{
			Summary:  item.Summary,
			Location: item.Location,
		}
		if item.Start.DateTime != "" {
			t, err := time.Parse(time.RFC3339, item.Start.DateTime)
			if err == nil {
				ce.StartTime = t.In(loc)
			}
		} else {
			ce.AllDay = true
			t, err := time.ParseInLocation("2006-01-02", item.Start.Date, loc)
			if err == nil {
				ce.StartTime = t
			}
		}
		if ce.Summary == "" {
			ce.Summary = "(タイトルなし)"
		}
		result = append(result, ce)
	}
	return result, nil
}
