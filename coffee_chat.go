package main

import (
	"log"

	misspell "github.com/client9/misspell"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	fuzzy "github.com/sajari/fuzzy"
)

const (
	helpStr = `FX Coffee bot helps team-mates buy beverages in the morning. It understands the following commands:

/coffee - every human places an order for hot beverage

/collect - one and only one human collects orders from the bot memory and physically places an aggregated order to buy and deliver beverages to the team-mates
After collection the bot is ready for the next round.

Bots can't initiate a conversation with human.`

	unknownCmdStr = `I'm sorry... I don't understand you...
Check available commands by typing /help.
If you need anything else, please speak to my manager @alelebe.`

	somethingWentWrongStr = `I'm sorry.. Something went wrong`
	pleaseTryAgainStr     = `Something went wrong... Please verify and try again`

	tellMeWhatToDoStr = "Please tell me what to do. Start with /"
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
	coffeeWatch   *CoffeeWatch
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
		p.replyToMessage(message, unknownCmdStr)
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
		p.replyToMessage(message, helpStr)

	case "/watch":
		if p.coffeeWatch == nil {
			p.coffeeWatch = initCoffeeWatch(p.Bot, message)
		}
		p.coffeeWatch.start()
	}
	return true
}

func (p *CoffeeChat) callbackQuery(callback tgbotapi.CallbackQuery) bool {

	if p.coffeeRequest != nil && p.coffeeRequest.isReplyOnMyMessage(callback) {
		p.coffeeRequest.onCallback(callback)
		return true
	}

	if p.coffeeCollect != nil && p.coffeeCollect.isReplyOnMyMessage(callback) {
		p.coffeeCollect.onCallback(callback)
		return true
	}

	if p.coffeeWatch != nil && p.coffeeWatch.isReplyOnMyMessage(callback) {
		p.coffeeWatch.onCallback(callback)
		return true
	}

	log.Printf("Coffee Chat: can't find msgId: %d in my chats, callback.Data: %s", callback.Message.MessageID, callback.Data)
	return false
}
