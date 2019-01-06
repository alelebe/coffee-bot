package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//CoffeeBot :
type CoffeeBot struct {
	Bot
	UpdateHandler

	activeChats map[int64]*CoffeeChat
}

func (p *CoffeeBot) getActiveChat(chatID int64) *CoffeeChat {
	if activeChat, ok := p.activeChats[chatID]; ok {
		return activeChat
	}
	return nil
}

//ProcessMessage : handling messages for the bot
func (p *CoffeeBot) ProcessMessage(message tgbotapi.Message) {

	chat := p.getActiveChat(message.Chat.ID)
	if chat == nil {
		chat = initCoffeeChat(p.Bot, message.Chat.ID)
		p.activeChats[message.Chat.ID] = chat
	}

	if message.Text[0] == '/' {
		chat.newCommand(message)
	} else {
		chat.newMessage(message)
	}
}

//ProcessCallback :
func (p *CoffeeBot) ProcessCallback(callback tgbotapi.CallbackQuery) {

	chat := p.getActiveChat(callback.Message.Chat.ID)
	if chat != nil && chat.callbackQuery(callback) {
		return
	}

	//Error handling
	if chat == nil {
		log.Printf("Coffee Bot: Can't find active chat by ID: %d", callback.Message.Chat.ID)
	}
	p.notifyUser(callback, "Please start again")
	// p.removeInlineKeyboard(callback)
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
		activeChats: map[int64]*CoffeeChat{},
	}

	bot.logBotDetails()
	return bot
}
