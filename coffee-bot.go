package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//CoffeeBot :
type CoffeeBot struct {
	Bot
	MessageHandler
}

//ProcessUpdate : handling messages for the bot
func (bot CoffeeBot) ProcessUpdate(update tgbotapi.Update) {

	bot.Send(tgbotapi.NewMessage(
		update.Message.Chat.ID,
		update.Message.Text,
	))
}

func initCoffeeBot(token string, debug bool) *CoffeeBot {
	if token == "" {
		return nil
	}

	bot := &CoffeeBot{
		Bot: Bot{
			"Coffee Bot",
			initBotAPI(token, debug),
		},
	}
	return bot
}
