package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Bot ...
type Bot struct {
	chatID int64
	api    *tgbotapi.BotAPI
}

// New ...
func New(token string, chatID int64, debug bool) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	api.Debug = debug
	return &Bot{
		api:    api,
		chatID: chatID,
	}, nil
}

// GetChatMembersCount ...
func (b *Bot) GetChatMembersCount() (membersCount int, err error) {
	membersCount, err = b.api.GetChatMembersCount(tgbotapi.ChatConfig{ChatID: b.chatID})
	if err != nil {
		log.Println(err)
	}
	return
}
