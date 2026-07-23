# Discord Calendar Reminder Bot

Googleカレンダーの「明日の予定」を、毎日決まった時刻にDiscordの指定チャンネルへ通知するBotです。

## 機能

- Google Calendar (OAuth2、自分のアカウント) から翌日の予定を取得
- 毎日指定した時刻に、常駐プロセス内のタイマーで自動実行
- Discordの指定チャンネルにEmbed形式で投稿

## 事前準備

### 1. Discord Botの作成

1. [Discord Developer Portal](https://discord.com/developers/applications) で新しいアプリケーションを作成
2. 「Bot」タブでBotを追加し、トークンを控える(`DISCORD_BOT_TOKEN`)
3. 「OAuth2 > URL Generator」で `bot` スコープと `Send Messages` `Embed Links` 権限を選び、生成されたURLからサーバーに招待
4. 投稿したいチャンネルを右クリック→「IDをコピー」(開発者モードを事前にON) → `DISCORD_CHANNEL_ID`

### 2. Google Calendar APIの有効化 & OAuth2クライアント作成

1. [Google Cloud Console](https://console.cloud.google.com/) でプロジェクトを作成
2. 「APIとサービス」→「ライブラリ」から **Google Calendar API** を有効化
3. 「APIとサービス」→「認証情報」→「認証情報を作成」→「OAuthクライアントID」
   - アプリケーションの種類は **デスクトップアプリ** を選択
4. 作成後、JSONをダウンロードし `credentials.json` としてプロジェクト直下に配置
   - OAuth同意画面でテストユーザーとして自分のGoogleアカウントを追加しておく(公開未申請の場合)

### 3. 環境変数の設定

```bash
cp .env.example .env
# .env を開いて DISCORD_BOT_TOKEN, DISCORD_CHANNEL_ID などを設定
```

## セットアップ & 実行

```bash
go mod tidy
go run .
```

初回起動時、ターミナルに認証用URLが表示されます。ブラウザで開いてGoogleアカウントでログイン・許可すると認可コードが表示されるので、それをターミナルに貼り付けてください。以後は `token.json` に保存され、自動でリフレッシュされます。

## ビルドして常駐させる例

```bash
go build -o calendar-reminder .
./calendar-reminder
```

Linuxサーバーで永続化する場合は `systemd` のserviceファイルにするか、`tmux`/`screen` 等で常駐させてください。

### systemd の例

```ini
[Unit]
Description=Discord Calendar Reminder Bot
After=network.target

[Service]
WorkingDirectory=/path/to/discord-calendar-reminder
ExecStart=/path/to/discord-calendar-reminder/calendar-reminder
Restart=always
EnvironmentFile=/path/to/discord-calendar-reminder/.env

[Install]
WantedBy=multi-user.target
```

## ディレクトリ構成

```
.
├── main.go       # エントリポイント・スケジューラ
├── config.go     # 環境変数の読み込み
├── calendar.go   # Google Calendar OAuth2認証・予定取得
├── discord.go    # Discordへの投稿処理
├── go.mod
├── .env.example
├── credentials.json
└── token.json
```

`.gitignore` に `credentials.json` `token.json` `.env` を追加

## カスタマイズ

- **複数カレンダーに対応したい場合**: `GOOGLE_CALENDAR_ID` をカンマ区切りにして `main.go` でループする形に拡張可能
- **DMにも送りたい場合**: `discord.go` に `session.UserChannelCreate(userID)` でDMチャンネルを作成し送信する処理を追加
- **当日の予定も通知したい場合**: `calendar.go` の `GetTomorrowEvents` を汎用化し、日付オフセットを引数化す