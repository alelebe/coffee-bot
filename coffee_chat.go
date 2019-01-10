package main

import (
	"log"

	misspell "github.com/client9/misspell"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	fuzzy "github.com/sajari/fuzzy"
)

//CoffeeChat :
type CoffeeChat struct {
	Bot
	chatID int64

	model    *fuzzy.Model
	replacer *misspell.StringReplacer

	// Command handlers
	coffeeRequest *CoffeeRequest
	coffeeCollect *CoffeeCollect
}

func initCoffeeChat(bot Bot, chatID int64) *CoffeeChat {
	oldnew := append(misspell.DictMain, misspell.DictBritish...)
	newChat := &CoffeeChat{
		Bot:      bot,
		chatID:   chatID,
		model:    fuzzy.NewModel(),
		replacer: misspell.NewStringReplacer(oldnew...),
	}
	return newChat
}

func (p *CoffeeChat) newMessage(message tgbotapi.Message) {
	text := p.replacer.Replace(message.Text)
	p.replyToMessage(message, text)
}

func (p *CoffeeChat) newCommand(message tgbotapi.Message) bool {

	switch message.Text {
	default:
		p.replyToMessage(message, "I'm sorry... I don't understand your command...")
		return false

	case "/coffee":
		// if p.coffeeCmd != nil {
		// 	p.coffeeCmd.cancel()
		// 	break
		// }
		if p.coffeeRequest == nil {
			p.coffeeRequest = initCoffeeRequest(p.Bot, message)
		}
		p.coffeeRequest.start()

	case "/collect":
		if p.coffeeCollect == nil {
			p.coffeeCollect = initCoffeeCollect(p.Bot, message)
		}
		p.coffeeCollect.start()

	case "/help":
		p.replyToMessage(message, "The bot helps team-mates buy hot beverages in the morning...\n\n"+
			"/coffee - every human places an order for hot beverage\n"+
			"/collect - one and only one human collects orders from the bot memory "+
			"and physically places an aggregated order to buy "+
			"and deliver beverages to the team-mates\n\n"+
			"After collection the bot is ready for the next round")
	}
	return true
}

func (p *CoffeeChat) callbackQuery(callback tgbotapi.CallbackQuery) bool {

	if p.coffeeRequest != nil {
		msg := p.coffeeRequest.isReplyOnMyMessage(callback)
		if msg != nil {
			p.coffeeRequest.onCallback(callback)
			return true
		}
	}

	if p.coffeeCollect != nil {
		msg := p.coffeeCollect.isReplyOnMyMessage(callback)
		if msg != nil {
			p.coffeeCollect.onCallback(callback)
			return true
		}
	}

	log.Printf("Old callback.Data '%s', skipping...", callback.Data)
	return false
}
