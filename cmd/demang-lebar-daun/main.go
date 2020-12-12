package main

import (
	"log"

	"github.com/palembang-digital/bots/pkg/configs"
	"github.com/palembang-digital/bots/pkg/googlesheet"
	"github.com/palembang-digital/bots/pkg/telegram"
)

func main() {
	var cfg = configs.Get()
	log.Println("Authenticating Telegram bot...")
	bot, err := telegram.New(cfg.TelegramToken, cfg.TelegramChatID, cfg.TelegramDebug)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Getting Telegram chat members count...")
	membersCount, err := bot.GetChatMembersCount()
	if err != nil {
		log.Panic(err)
	}
	log.Println("Telegram chat members count:", membersCount)

	log.Println("Authenticating to Google Sheets...")
	sheetsService, err := googlesheet.New(cfg.SheetsCredentials, cfg.SheetsSpreadsheetID, cfg.SheetsSpreadsheetRange)
	if err != nil {
		log.Panic(err)
	}
	err = sheetsService.AppendMembersCount(membersCount)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Job completed! See ya~")
}
