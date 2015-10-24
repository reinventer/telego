# Telego

A simple wrapper for create telegram bot

## Install

```sh
$ go get github.com/reinventer/telego
```

## Make Bot

You need to talk with [BotFather](https://telegram.me/botfather) and follow a few simple steps. When you've created a bot, you received your authorization token. Save it.

Write a simple code:

```go
package main

import "github.com/reinventer/telego"

func main() {
	b, err := telego.NewBot("token")
	if err != nil {
		panic(err)
	}

	b.SetHandlerWithHelp("/test", "show test message", TestHandler)
	b.SetDefaultHandler(DefaultHandler)
	b.Run()
}

func TestHandler(update *telego.Update) []string {
	return []string{
		"It's a test message, " + update.Message.From.UserName,
		"Parameter: '" + update.Params + "'",
	}
}

func DefaultHandler(update *telego.Update) []string {
	return []string{"Unknown comand", "Try to use /help"}
}
```

Compile and run it.

Now talk with your bot, try commands `/test`, `/test params`, `/help`

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).
