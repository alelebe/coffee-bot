package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//CoffeeBot :
type CoffeeBot struct {
	Bot
	UpdateHandler

	activeChats  map[int64]*CoffeeChat
	lostMessages map[int]interface{}
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
	if chat == nil {

		log.Printf("Coffee Bot: Can't find active chat by ID: %d", callback.Message.Chat.ID)

	} else if chat.callbackQuery(callback) {
		return
	}

	p.removeLostMessages(callback)
}

func (p *CoffeeBot) removeLostMessages(callback tgbotapi.CallbackQuery) {
	ID := callback.Message.MessageID

	if _, ok := p.lostMessages[ID]; ok { //found

		p.removeInlineKeyboard(callback)
		delete(p.lostMessages, ID)

	} else {

		p.alertUser(callback, "The history was lost...\nPlease click on the same item again\nto remove that inline keyboard from the chat")
		p.lostMessages[ID] = true
	}
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
		activeChats:  make(map[int64]*CoffeeChat, 0),
		lostMessages: make(map[int]interface{}, 0),
	}

	bot.logBotDetails()
	return bot
}
