package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func NewCalenderService(ctx context.Context, cfg *Config) (*calender.Service, error) {
	b, err := os.ReadFile(cfg.CredentialsPath)
	if err != nil {
		return nil, fmt.Errorf("credentials.json の読み込みに失敗しました (%s): %w", cfg.CredentialsPath, err)
	}

	oauthConfig, err := google.ConfigFromJSON(b, calender.CalendarReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("credentials.json のパースに失敗しました: %w", err)
	}

	client, err := getHTTPClient(oauthConfig, cfg.TokenPath)
	if err != nil {
		return nil, err
	}

	srv, err := calender.NewService(ctx, option.WithHTTPClient(client))
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

	ts := config.TokenSource(config.Background(), tok)
	return oauth2.NewClient(config.Background(), tokenSever{ts, tokenPath}), nil
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Println("=== Google Calendar 初回認証 ===")
	fmt.Println("以下のURLをブラウザで開き、Googleアカウントで許可してください:")
	fmt.Println(authURL)
	fmt.Print("表示された認可コードを貼り付けてEnterを押してください: ")

	var authCode string
	
}
