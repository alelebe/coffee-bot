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

func initCoffeeCollect(bot Bot, message tgbotapi.Message) *CoffeeCollect {
	newCmd := &CoffeeCollect{
		Bot:        bot,
		initialMsg: message,
		chatID:     message.Chat.ID,

		myRequests: make(map[int]CollectionRequest, 0),
	}
	return newCmd
}

func (p *CoffeeCollect) start() {

	orders, cas := ordersReadyForCollection()
	if len(orders) == 0 {
		p.replyToMessage(p.initialMsg, "No orders are ready for collection...\nPlease try again later")
		return
	}

	sent, err := p.replyToMessageWithInlineKeyboard(
		p.initialMsg,
		ordersForCollection(orders),
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

func ordersForCollection(orders []CoffeeOrder) string {
	var s string
	if len(orders) > 1 {
		s = "s"
	}

	result := fmt.Sprintf("*%d* order%s ready for collection:\n", len(orders), s)
	for _, it := range orders {
		result += "\n*" + it.UserName + "*: " + it.Beverage
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

func (p CoffeeCollect) isReplyOnMyMessage(callback tgbotapi.CallbackQuery) bool {
	if callback.Message != nil {
		if _, ok := p.myRequests[callback.Message.MessageID]; ok {
			return true
		}
	}
	return false
}

func (p *CoffeeCollect) onCallback(callback tgbotapi.CallbackQuery) {
	//have to resolve msgID to internal request (again)!
	req, ok := p.myRequests[callback.Message.MessageID]
	if !ok {
		return
	}

	button := callback.Data

	switch button {
	default:
		return

	case btnCOLLECT:
		p.finishRequest(callback, req)

	case btnCANCEL:
		p.updateMessage(callback, "Collection cancelled...")
		p.removeInlineKeyboard(callback)
		// p.notifyUser(callback, "Request aborted")
	}

	//remove request from my queue
	delete(p.myRequests, callback.Message.MessageID)
}

func (p *CoffeeCollect) finishRequest(callback tgbotapi.CallbackQuery, request CollectionRequest) {

	if collectOrdes(request.cas) {
		log.Printf("Coffee Collect: request is successfully collected: %+v", request)
		p.updateMessage(callback, collectedOrders(request.orders))
		p.notifyOnCollection(callback.Message.Chat.ID, request.orders)

	} else {
		p.updateMessage(callback, pleaseTryAgainStr)
	}
	p.removeInlineKeyboard(callback)
}

func collectedOrders(orders []CoffeeOrder) string {

	beverages := make(map[string]int, 0)
	for _, it := range orders {
		beverages[it.Beverage]++
	}

	result := "You've just collected:"
	for it, value := range beverages {
		result += fmt.Sprintf("\n*%d*\t%s", value, it)
	}
	return result
}

func (p *CoffeeCollect) notifyOnCollection(originalChatID int64, orders []CoffeeOrder) {

	//remember User IDs of chats where the notification was sent
	alreadyNotified := make(map[int]bool, 0)
	alreadyNotified[p.initialMsg.From.ID] = true

	//notify those places orders for coffee
	yourOrderWasCollected := fmt.Sprintf("Your order was collected by %s", p.initialMsg.From.FirstName)
	for _, obj := range orders {
		if _, ok := alreadyNotified[obj.UserID]; ok {
			continue
		}
		p.sendToChat(obj.ChatID, obj.UserName+",\n"+yourOrderWasCollected)
		alreadyNotified[obj.UserID] = true
	}

	//notify all coffee watchers
	orderWasCollected := fmt.Sprintf("The order was collected by\n*%s*", p.initialMsg.From.FirstName)
	for _, obj := range allCoffeeWatchers() {
		if _, ok := alreadyNotified[obj.UserID]; ok {
			continue
		}
		p.sendToChat(obj.ChatID, orderWasCollected)
		alreadyNotified[obj.UserID] = true
	}
}
