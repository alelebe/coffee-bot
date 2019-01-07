package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// CoffeeCollect :
type CoffeeCollect struct {
	Bot
	initialMsg tgbotapi.Message
	chatID     int64

	// orders []CoffeeOrder
}

func (p *CoffeeCollect) start() {

	orders := collectOrders()
	log.Printf("%d orders ready for collection: %+v", len(orders), orders)

	if len(orders) == 0 {
		p.replyToMessage(p.initialMsg, "No orders ready for collection :(")
		return
	}
	_, err := p.replyToMessageWithInlineKeyboard(
		p.initialMsg,
		fmt.Sprintf("%d oders from:", len(orders))+ordersFrom(orders),
		confirmCollection(),
	)
	if err == nil {
		// p.myMessages[sent.MessageID] = sent
	}
}

func ordersFrom(orders []CoffeeOrder) string {
	var result string
	for _, it := range orders {
		result += it.UserName
	}
	return result
}

func confirmCollection() [][]tgbotapi.InlineKeyboardButton {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(btnCOLLECT, btnCOLLECT),
		tgbotapi.NewInlineKeyboardButtonData(btnCANCEL, btnCANCEL),
	))
	return keyboard
}

func initCoffeeCollect(bot Bot, message tgbotapi.Message) *CoffeeCollect {
	newCmd := &CoffeeCollect{
		Bot:        bot,
		initialMsg: message,
		chatID:     message.Chat.ID,
	}
	return newCmd
}
