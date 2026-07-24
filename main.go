package main

import (
	"context"
	"log"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	ctx := context.Background()
}
