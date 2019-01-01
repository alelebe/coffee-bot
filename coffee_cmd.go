package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//CoffeeCmd :
type CoffeeCmd struct {
	Bot
	initialMsg tgbotapi.Message
	chatID     int64

	lastMessage *tgbotapi.Message
}

func (p CoffeeCmd) start() {
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Cappuchino", "Cappuchino"),
			tgbotapi.NewInlineKeyboardButtonData("Americano", "Americano"),
		))

	msg := tgbotapi.NewMessage(
		p.chatID,
		"Choose a drink from the list below:",
	)
	// msg.BaseChat.ReplyToMessageID = original.MessageID
	msg.BaseChat.ReplyMarkup = buttons
	p.Send(msg)
}

func (p *CoffeeCmd) cancel() {
	if p.lastMessage != nil {
		p.Send(
			tgbotapi.NewDeleteMessage(
				p.chatID,
				p.lastMessage.MessageID,
			))
		p.lastMessage = nil
	}
}

func isSameDate(date1 int, date2 int) bool {
	y1, m1, d1 := time.Unix(int64(date1), 0).Date()
	y2, m2, d2 := time.Unix(int64(date2), 0).Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (p *CoffeeCmd) nextStep(callback tgbotapi.CallbackQuery) {
	if p.lastMessage == nil {
		log.Printf("CoffeeCmd: Last message is empty\n")
		return
	}
	if p.lastMessage.MessageID != callback.Message.MessageID {
		log.Printf("CoffeeCmd: Callback is called on wrong message\n")
		if isSameDate(p.lastMessage.Date, callback.Message.Date) {
			p.Send(
				tgbotapi.NewDeleteMessage(
					p.chatID,
					callback.Message.MessageID,
				))
		}
		return
	}

	p.AnswerCallbackQuery(
		tgbotapi.NewCallback(
			callback.ID,
			fmt.Sprintf("%s\nGood choice!", callback.Data),
		))

	msg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "Would you like some chocolate on top?")
	_, err := p.Send(msg)
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
	_, err = p.Send(msg2)
	if err != nil {
		log.Printf("Failed to edit message: %v\n", err)
	}
}

func initCoffeeCmd(bot Bot, message tgbotapi.Message) *CoffeeCmd {
	newCmd := &CoffeeCmd{
		Bot:        bot,
		initialMsg: message,
		chatID:     message.Chat.ID,
	}
	return newCmd
}
