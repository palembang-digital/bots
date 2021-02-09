package configs

import (
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

// Config stores configuration for the bot.
type Config struct {
	TelegramToken  string `required:"true" split_words:"true"`
	TelegramChatID int64  `required:"true" split_words:"true"`
	TelegramDebug  bool   `split_words:"true"`

	SheetsCredentials      string `split_words:"true"`
	SheetsSpreadsheetID    string `split_words:"true"`
	SheetsSpreadsheetRange string `split_words:"true" default:"Telegram Bot"`
}

var conf Config
var once sync.Once

// Get returns the singleton config instance
func Get() Config {
	once.Do(func() {
		err := envconfig.Process("", &conf)
		if err != nil {
			log.Fatal("Can't load config: ", err)
		}
	})
	return conf
}
