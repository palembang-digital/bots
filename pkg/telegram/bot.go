package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Bot object.
type Bot struct {
	chatID int64
	api    *tgbotapi.BotAPI
}

// New initializes a new Bot object.
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

// GetChatMembersCount gets the number of members in given chat group.
func (b *Bot) GetChatMembersCount() (membersCount int, err error) {
	membersCount, err = b.api.GetChatMembersCount(tgbotapi.ChatConfig{ChatID: b.chatID})
	if err != nil {
		log.Println(err)
	}
	return
}

// Send sends the text message to given chat ID.
// It sends the message in HTML parse mode.
func (b *Bot) Send(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML

	m, err := b.api.Send(msg)
	if err != nil {
		return err
	}
	log.Printf("%+v", m)
	return nil
}
