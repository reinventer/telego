package telego

import (
	"strings"

	"github.com/Syfaro/telegram-bot-api"
)

type Bot struct {
	Api      *tgbotapi.BotAPI
	handlers map[string]func(string, tgbotapi.Update)
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{Api: api}, nil
}

func (b *Bot) AddHandler(command string, handler func(string, tgbotapi.Update)) {
	b.handlers[command] = handler
}

func (b *Bot) Run() {
	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	err := b.Api.UpdatesChan(ucfg)
	// читаем обновления из канала
	for {
		select {
		case update := <-b.Api.Updates:
			text := strings.SplitN(update.Message.Text, ' ', 2)
			handler, ok := b.handlers[text[0]]
			if ok {
				handler(text[1], update)
			}
		}
	}
}
