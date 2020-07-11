package main

import (
	"context"
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Config stores configuration for the bot.
type Config struct {
	TelegramToken  string `required:"true" split_words:"true"`
	TelegramChatID int64  `required:"true" split_words:"true"`
	TelegramDebug  bool   `split_words:"true"`

	SheetsCredentials      string `required:"true" split_words:"true"`
	SheetsSpreadsheetID    string `required:"true" split_words:"true"`
	SheetsSpreadsheetRange string `split_words:"true" default:"Telegram Bot"`
}

func main() {
	ctx := context.Background()

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Panic(err)
	}

	log.Println("Authenticating Telegram bot...")
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = cfg.TelegramDebug

	log.Println("Getting Telegram chat members count...")
	membersCount, err := bot.GetChatMembersCount(tgbotapi.ChatConfig{ChatID: cfg.TelegramChatID})
	if err != nil {
		log.Panic(err)
	}

	log.Println("Telegram chat members count:", membersCount)

	log.Println("Authenticating to Google Sheets...")
	sheetsService, err := sheets.NewService(ctx, option.WithCredentialsJSON([]byte(cfg.SheetsCredentials)))
	if err != nil {
		log.Panic(err)
	}

	log.Println("Writing to Google Sheets...")
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{{time.Now().UTC(), time.Now().Month(), membersCount}},
	}
	resp, err := sheetsService.Spreadsheets.Values.
		Append(cfg.SheetsSpreadsheetID, cfg.SheetsSpreadsheetRange, valueRange).
		ValueInputOption("RAW").
		Do()
	if err != nil {
		log.Panic(err)
	}

	if resp.ServerResponse.HTTPStatusCode != http.StatusOK {
		log.Panic(resp.ServerResponse.HTTPStatusCode)
	}

	log.Println("Job completed! See ya~")
}
