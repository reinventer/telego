package telego

import (
	"strings"

	"github.com/Syfaro/telegram-bot-api"
)

type Bot struct {
	Api      *tgbotapi.BotAPI
	handlers map[string]func(*Update) string
}

type Update struct {
	tgbotapi.Update
	Params string
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{Api: api, handlers: map[string]func(*Update) string{}}, nil
}

func (b *Bot) AddHandler(command string, handler func(*Update) string) {
	b.handlers[command] = handler
}

func (b *Bot) Run() {
	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	b.Api.UpdatesChan(ucfg)
	for {
		select {
		case tupdate := <-b.Api.Updates:
			text := strings.SplitN(tupdate.Message.Text, " ", 2)
			handler, ok := b.handlers[text[0]]
			if ok {
				params := ""
				if len(text) > 1 {
					params = text[1]
				}
				update := Update{
					Update: tupdate,
					Params: params,
				}
				msg_text := handler(&update)
				if msg_text != "" {
					chat_id := tupdate.Message.Chat.ID
					msg := tgbotapi.NewMessage(chat_id, msg_text)
					b.Api.SendMessage(msg)
				}
			}
		}
	}
}
