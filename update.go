package telego

import "github.com/Syfaro/telegram-bot-api"

type Update struct {
    tgbotapi.Update
    Bot    *Bot
    Params string
}

func (u *Update)Reply(text string) error {
    return u.Bot.SendTextMessage(u.Message.Chat.ID, text)
}
