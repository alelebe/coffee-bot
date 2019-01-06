package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// CoffeeCollect :
type CoffeeCollect struct {
	Bot
	initialMsg tgbotapi.Message
	chatID     int64
	// myMessages map[int]tgbotapi.Message
}

func (p *CoffeeCollect) start() {

	orders := collectOrders()
	log.Printf("%d orders ready for collection: %+v", len(orders), orders)
}

func initCoffeeCollect(bot Bot, message tgbotapi.Message) *CoffeeCollect {
	newCmd := &CoffeeCollect{
		Bot:        bot,
		initialMsg: message,
		chatID:     message.Chat.ID,
		// myMessages: make(map[int]tgbotapi.Message),
	}
	return newCmd
}
