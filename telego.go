package telego

import (
	"strings"
	"sync"

	"github.com/Syfaro/telegram-bot-api"
)

type Bot struct {
	Api             *tgbotapi.BotAPI
	handlers        map[string]func(*Update) []string
	default_handler func(*Update) []string
}

type Update struct {
	tgbotapi.Update
	Bot    *tgbotapi.BotAPI
	Params string
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot := Bot{
		Api:      api,
		handlers: map[string]func(*Update) []string{},
	}
	return &bot, nil
}

func (b *Bot) AddHandler(command string, handler func(*Update) []string) {
	b.handlers[command] = handler
}

func (b *Bot) DefaultHandler(handler func(*Update) []string) {
	b.default_handler = handler
}

func (b *Bot) RunWorkers(workers_count int) {
	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	b.Api.UpdatesChan(ucfg)
	if workers_count < 1 {
		workers_count = 1
	}
	var wg sync.WaitGroup
	wg.Add(workers_count)
	for i := 0; i < workers_count; i++ {
		wg.Add(1)
		go b.worker(b.Api.Updates, &wg)
	}
	wg.Wait()
}

func (b *Bot) Run() {
	b.RunWorkers(1)
}

func (b *Bot) newUpdate(tupdate tgbotapi.Update, params string) *Update {
	return &Update{
		Update: tupdate,
		Bot:    b.Api,
		Params: params,
	}
}

func (b *Bot) worker(updates <-chan tgbotapi.Update, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case tupdate := <-updates:
			text := strings.SplitN(tupdate.Message.Text, " ", 2)
			handler, ok := b.handlers[text[0]]
			if ok {
				params := ""
				if len(text) > 1 {
					params = text[1]
				}

				b.execHandler(handler, b.newUpdate(tupdate, params))
			} else if b.default_handler != nil {
				b.execHandler(b.default_handler, b.newUpdate(tupdate, ""))
			}
		}
	}
}

func (b *Bot) execHandler(handler func(*Update) []string, update *Update) {
	texts := handler(update)
	for _, msg_text := range texts {
		if msg_text != "" {
			chat_id := update.Message.Chat.ID
			msg := tgbotapi.NewMessage(chat_id, msg_text)
			b.Api.SendMessage(msg)
		}
	}
}
