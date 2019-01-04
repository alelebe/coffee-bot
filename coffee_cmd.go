package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type node struct {
	item   *Drink
	parent *Drink
}

//CoffeeCmd :
type CoffeeCmd struct {
	Bot
	initialMsg tgbotapi.Message
	chatID     int64

	drinks      []Drink
	allDrinks   map[string]node
	lastMessage *tgbotapi.Message
}

func buildKeyboard(drinks []Drink) [][]tgbotapi.InlineKeyboardButton {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	numOfRows := len(drinks) / 2
	row := 0
	for idx := 0; idx < len(drinks); idx++ {
		item := drinks[idx]

		if row < numOfRows {
			nextItem := drinks[idx+1]
			keys := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(item.Display, item.Display),
				tgbotapi.NewInlineKeyboardButtonData(nextItem.Display, nextItem.Display),
			)
			idx++
			keyboard = append(keyboard, keys)
		} else {
			keys := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(item.Display, item.Display),
			)
			keyboard = append(keyboard, keys)
		}
	}
	return keyboard
}

func (p *CoffeeCmd) start() {

	msg := tgbotapi.NewMessage(
		p.chatID,
		"Choose a drink from the list below:",
	)
	msg.BaseChat.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buildKeyboard(p.drinks)...)
	sent, _ := p.Send(msg)
	p.lastMessage = &sent
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

func buildTree(parent *Drink, drinks []Drink, tree map[string]node) {
	for _, item := range drinks {
		tree[item.Display] = node{
			item:   &item,
			parent: parent,
		}
		buildTree(&item, item.SubItems, tree)
	}
}

func initCoffeeCmd(bot Bot, message tgbotapi.Message) *CoffeeCmd {
	newCmd := &CoffeeCmd{
		Bot:        bot,
		initialMsg: message,
		chatID:     message.Chat.ID,
	}
	const filePath = "./data/benugo.json"
	menu, err := loadBewerages(filePath)
	if err != nil {
		dir, _ := os.Getwd()
		log.Printf("%s", dir)
		log.Fatal(err)
	}
	newCmd.drinks = menu.Bewerages
	newCmd.allDrinks = make(map[string]node, 0)
	buildTree(nil, newCmd.drinks, newCmd.allDrinks)
	log.Printf("Loaded hot bewerages from %s", filePath)
	log.Printf("Number of available items: %d\n", len(newCmd.allDrinks))
	return newCmd
}
