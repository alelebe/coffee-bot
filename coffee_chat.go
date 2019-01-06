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

	coffeeCmd *CoffeeCmd
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
	case "/coffee":
		// if p.coffeeCmd != nil {
		// 	p.coffeeCmd.cancel()
		// 	break
		// }
		if p.coffeeCmd == nil {
			p.coffeeCmd = initCoffeeCmd(p.Bot, message)
		}
		p.coffeeCmd.start()

	case "/collect":
	default:
		p.replyToMessage(message, "I'm sorry... I don't understand your command...")
		return false
	}
	return true
}

func (p *CoffeeChat) callbackQuery(callback tgbotapi.CallbackQuery) bool {

	if p.coffeeCmd != nil {
		msg := p.coffeeCmd.isReplyOnMyMessage(callback)
		if msg != nil {
			p.coffeeCmd.onCallback(callback, *msg)
			return true
		}
	}
	log.Printf("Old callback.Data '%s', skipping...", callback.Data)
	return false
}
