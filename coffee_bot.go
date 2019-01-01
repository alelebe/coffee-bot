package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//CoffeeBot :
type CoffeeBot struct {
	Bot
	UpdateHandler

	activeChats map[int64]CoffeeChat
}

func (bot CoffeeBot) getActiveChat(chatID int64) CoffeeChat {
	if activeChat, ok := bot.activeChats[chatID]; ok {
		return activeChat
	}

	newChat := initCoffeeChat(bot.Bot, chatID)
	bot.activeChats[chatID] = newChat
	return newChat
}

//ProcessMessage : handling messages for the bot
func (bot CoffeeBot) ProcessMessage(message tgbotapi.Message) {

	chat := bot.getActiveChat(message.Chat.ID)

	if message.Text[0] == '/' {
		go chat.newCommand(message)
	} else {
		go chat.newMessage(message)
	}
}

//ProcessCallback :
func (bot CoffeeBot) ProcessCallback(callback tgbotapi.CallbackQuery) {

	chat := bot.getActiveChat(callback.Message.Chat.ID)
	chat.callbackQuery(callback)
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
		activeChats: map[int64]CoffeeChat{},
	}

	bot.logBotDetails()
	return bot
}
