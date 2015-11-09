package telego

import (
	"strings"
	"sync"

	"github.com/Syfaro/telegram-bot-api"
)

type HandlerFunc func(*Update)

type Bot struct {
	sync.RWMutex
	Api                   *tgbotapi.BotAPI
	handlers              map[string]HandlerFunc
	default_handler       HandlerFunc
	handlers_descriptions map[string]string
	handlers_order        []string
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot := Bot{
		Api:                   api,
		handlers:              map[string]HandlerFunc{},
		handlers_descriptions: map[string]string{},
		handlers_order:        []string{},
	}
	return &bot, nil
}

func (b *Bot) SetHandler(command string, handler HandlerFunc) {
	b.Lock()
	b.handlers[command] = handler
	b.Unlock()
}

func (b *Bot) SetHandlerWithHelp(command string, description string, handler HandlerFunc) {
	b.Lock()
	b.handlers[command] = handler
	b.handlers_descriptions[command] = description
	for i, cmd := range b.handlers_order {
		if cmd == command {
			b.handlers_order = append(b.handlers_order[:i], b.handlers_order[i+1:]...)
			break
		}
	}
	b.handlers_order = append(b.handlers_order, command)
	b.handlers["/help"] = b.defaultHelpHandler
	b.handlers["/start"] = b.defaultHelpHandler
	b.Unlock()
}

func (b *Bot) SetDefaultHandler(handler HandlerFunc) {
	b.default_handler = handler
}

func (b *Bot) RunWorkers(workers_count int) {
	b.initUpdates()
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
	b.initUpdates()
	for tupdate := range b.Api.Updates {
		go b.findAndExecHandler(tupdate)
	}
}

func (b *Bot) defaultHelpHandler(update *Update) {
	help_message := ""
	b.RLock()
	for _, cmd := range b.handlers_order {
		help_message += cmd + " - " + b.handlers_descriptions[cmd] + "\n"
	}
	b.RUnlock()
	update.Reply(help_message)
}

func (b *Bot) newUpdate(tupdate tgbotapi.Update, params string) *Update {
	return &Update{
		Update: tupdate,
		Bot:    b,
		Params: params,
	}
}

func (b *Bot) SendTextMessage(chat_id int, text string) error {
	msg := tgbotapi.NewMessage(chat_id, text)
	_, err := b.Api.SendMessage(msg)
	return err
}

func (b *Bot) initUpdates() {
	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	b.Api.UpdatesChan(ucfg)
}

func (b *Bot) worker(updates <-chan tgbotapi.Update, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case tupdate := <-updates:
			b.findAndExecHandler(tupdate)
		}
	}
}

func (b *Bot) findAndExecHandler(tupdate tgbotapi.Update) {
	text := strings.SplitN(tupdate.Message.Text, " ", 2)
	handler, ok := b.handlers[text[0]]
	if ok {
		params := ""
		if len(text) > 1 {
			params = text[1]
		}

		handler(b.newUpdate(tupdate, params))
	} else if b.default_handler != nil {
		b.default_handler(b.newUpdate(tupdate, ""))
	}
}
