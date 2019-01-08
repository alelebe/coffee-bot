package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// CollectionRequest :
type CollectionRequest struct {
	message    tgbotapi.Message
	orders     []CoffeeOrder
	orders_cas uint64
}

// CoffeeCollect :
type CoffeeCollect struct {
	Bot
	initialMsg tgbotapi.Message
	chatID     int64

	myRequests map[int]CollectionRequest
}

func (p *CoffeeCollect) start() {

	var orders []CoffeeOrder
	var cas uint64

	orders, cas = collectOrders()
	log.Printf("%d orders ready for collection: %+v", len(orders), orders)

	if len(orders) == 0 {
		p.replyToMessage(p.initialMsg, "No orders are ready for collection today...\nPlease try again later")
		return
	}

	sent, err := p.replyToMessageWithInlineKeyboard(
		p.initialMsg,
		ordersFromUsers(orders),
		confirmCollection(),
	)
	if err == nil {
		p.myRequests[sent.MessageID] = CollectionRequest{
			message:    sent,
			orders:     orders,
			orders_cas: cas,
		}
	}
}

func ordersFromUsers(orders []CoffeeOrder) string {
	result := fmt.Sprintf("I have %d orders ready for collection from those users:\n", len(orders))
	for _, it := range orders {
		result += "\n*" + it.UserName + "*"
	}
	if len(orders) > 0 {
		result += "\n\nPlease confirm you would like to collect?"
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

func (p CoffeeCollect) isReplyOnMyMessage(callback tgbotapi.CallbackQuery) *tgbotapi.Message {
	if callback.Message != nil {
		if req, ok := p.myRequests[callback.Message.MessageID]; ok {
			return &req.message
		}
	}
	log.Printf("Coffee Collect: can't find msgId: %d in my messages: %+v", callback.Message.MessageID, p.myRequests)
	return nil
}

func (p *CoffeeCollect) onCallback(callback tgbotapi.CallbackQuery) {
	//have to resolve msgID to internal request (again)!
	req, ok := p.myRequests[callback.Message.MessageID]
	if !ok {
		return
	}

	button := callback.Data

	switch button {

	case btnCONFIRM:
		p.finishRequest(callback, req)
	}
}

func (p *CoffeeCollect) finishRequest(callback tgbotapi.CallbackQuery, request CollectionRequest) {

	log.Printf("Coffee Collect: drinks are ready for collection: %+v", request)

}

func initCoffeeCollect(bot Bot, message tgbotapi.Message) *CoffeeCollect {
	newCmd := &CoffeeCollect{
		Bot:        bot,
		initialMsg: message,
		chatID:     message.Chat.ID,

		myRequests: make(map[int]CollectionRequest, 0),
	}
	return newCmd
}
