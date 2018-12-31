package main

import (
	"fmt"
	"log"

	misspell "github.com/client9/misspell"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	fuzzy "github.com/sajari/fuzzy"
)

//CoffeeChat :
type CoffeeChat struct {
	bot  Bot
	chat *tgbotapi.Chat

	model    *fuzzy.Model
	replacer *misspell.StringReplacer
}

func initCoffeeChat(bot Bot, chat *tgbotapi.Chat) CoffeeChat {
	oldnew := append(misspell.DictMain, misspell.DictBritish...)
	newChat := CoffeeChat{
		bot:      bot,
		chat:     chat,
		model:    fuzzy.NewModel(),
		replacer: misspell.NewStringReplacer(oldnew...),
	}
	return newChat
}

func (p *CoffeeChat) replyTo(original tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(
		p.chat.ID,
		text,
	)
	//msg.BaseChat.ReplyToMessageID = original.MessageID
	p.bot.Send(msg)
}
func (p *CoffeeChat) replyToWithMarkup(original tgbotapi.Message, text string, markup tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(
		p.chat.ID,
		text,
	)
	//msg.BaseChat.ReplyToMessageID = original.MessageID
	msg.BaseChat.ReplyMarkup = markup
	p.bot.Send(msg)
}

func (p *CoffeeChat) newMessage(message tgbotapi.Message) {
	text := p.replacer.Replace(message.Text)
	p.replyTo(message, text)
	// tgbotapi.ReplyKeyboardMarkup()
}

func (p *CoffeeChat) newCommand(message tgbotapi.Message) bool {

	switch message.Text {
	case "/coffee":
		buttons := p.newCoffeeBoard()
		p.replyToWithMarkup(message, "Choose a drink from the list below:", buttons)

	case "/collect":
	default:
		return false
	}
	return true
}

func (p CoffeeChat) newCoffeeBoard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Cappuchino", "Cappuchino"),
			tgbotapi.NewInlineKeyboardButtonData("Americano", "Americano"),
		))
}

func (p CoffeeChat) callbackQuery(callback tgbotapi.CallbackQuery) {

	p.bot.AnswerCallbackQuery(
		tgbotapi.NewCallback(
			callback.ID,
			fmt.Sprintf("%s\nGood choice!", callback.Data),
		))

	msg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "Would you like some chocolate on top?")
	_, err := p.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to edit message: %v\n", err)
	}
	msg2 := tgbotapi.NewEditMessageReplyMarkup(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Yes, please", "1"),
				tgbotapi.NewInlineKeyboardButtonData("No, thanks", "0"),
			)),
	)
	_, err = p.bot.Send(msg2)
	if err != nil {
		log.Printf("Failed to edit message: %v\n", err)
	}
}
