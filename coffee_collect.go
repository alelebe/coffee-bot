package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// CollectionRequest :
type CollectionRequest struct {
	message tgbotapi.Message
	orders  []CoffeeOrder
	cas     uint64
}

// CoffeeCollect :
type CoffeeCollect struct {
	Bot
	initialMsg tgbotapi.Message
	chatID     int64

	myRequests map[int]CollectionRequest
}

func (p *CoffeeCollect) start() {

	orders, cas := ordersReadyForCollection()
	if len(orders) == 0 {
		p.replyToMessage(p.initialMsg, "No orders are ready for collection today...\nPlease try again later")
		return
	}

	sent, err := p.replyToMessageWithInlineKeyboard(
		p.initialMsg,
		ordersFromUsers(orders),
		confirmCollection(),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d order(s) ready for collection: %+v", len(orders), orders)

	p.myRequests[sent.MessageID] = CollectionRequest{
		message: sent,
		orders:  orders,
		cas:     cas,
	}
}

func ordersFromUsers(orders []CoffeeOrder) string {
	var s string
	if len(orders) > 1 {
		s = "s"
	}

	result := fmt.Sprintf("*%d* order%s ready for collection from:\n", len(orders), s)
	for _, it := range orders {
		result += "\n*" + it.UserName + "*: " + it.Beverage
	}

	if len(orders) > 0 {
		result += "\n\nPlease confirm you would like to collect?"
	}
	return result
}

func ordersToCollect(orders []CoffeeOrder) string {

	bewerages := make(map[string]int, 0)
	for _, it := range orders {
		bewerages[it.Beverage]++
	}

	result := "You've just collected:"
	for it, value := range bewerages {
		result += fmt.Sprintf("\n*%d*\t%s", value, it)
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
	case btnCOLLECT:
		p.finishRequest(callback, req)

	case btnCANCEL:
		p.updateMessage(callback, "Collection cancelled...")
		p.removeInlineKeyboard(callback)
		// p.notifyUser(callback, "Request aborted")
	}
}

func (p *CoffeeCollect) finishRequest(callback tgbotapi.CallbackQuery, request CollectionRequest) {

	if collectOrdes(request.cas) {
		log.Printf("Coffee Collect: request is successfully collected: %+v", request)
		p.updateMessage(callback, ordersToCollect(request.orders))

	} else {
		p.updateMessage(callback, "New order has just arrived... please verify and start again...")
	}
	p.removeInlineKeyboard(callback)
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
