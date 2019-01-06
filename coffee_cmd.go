package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	btnOK     = "OK"
	btnCANCEL = "Cancel"
)

// type node struct {
// 	item   *Drink
// 	parent *Drink
// }

//CoffeeCmd :
type CoffeeCmd struct {
	Bot
	initialMsg tgbotapi.Message
	chatID     int64

	entry Bewerages
	// allDrinks  map[string]node
	myMessages map[int]tgbotapi.Message
}

func chooseOneDrinkKeyboard(drinks []Drink, hasBackButton bool) [][]tgbotapi.InlineKeyboardButton {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	numOfRows := len(drinks) / 2
	row := 0
	for idx := 0; idx < len(drinks); idx++ {
		item := drinks[idx]

		if row < numOfRows {
			nextItem := drinks[idx+1]
			keys := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(item.Display, item.ID),
				tgbotapi.NewInlineKeyboardButtonData(nextItem.Display, nextItem.ID),
			)
			idx++
			keyboard = append(keyboard, keys)
		} else {
			keys := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(item.Display, item.ID),
			)
			keyboard = append(keyboard, keys)
		}
	}
	return keyboard
}

func confirmDrinkKeyboard(item Drink) [][]tgbotapi.InlineKeyboardButton {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(btnOK, fmt.Sprintf("%s::%s", btnOK, item.ID)),
		tgbotapi.NewInlineKeyboardButtonData(btnCANCEL, fmt.Sprintf("%s::%s", btnCANCEL, item.ID)),
	))
	return keyboard
}

func (p *CoffeeCmd) start() {

	sent, err := p.replyToMessageWithInlineKeyboard(
		p.initialMsg, p.entry.Question,
		chooseOneDrinkKeyboard(p.entry.Items, false),
	)
	if err == nil {
		p.myMessages[sent.MessageID] = sent
	}
}

// func (p *CoffeeCmd) cancel() {
// 	if p.lastMessage != nil {
// 		p.Send(
// 			tgbotapi.NewDeleteMessage(
// 				p.chatID,
// 				p.lastMessage.MessageID,
// 		))
// 		p.lastMessage = nil
// 	}
// }

// func isSameDate(date1 int, date2 int) bool {
// 	y1, m1, d1 := time.Unix(int64(date1), 0).Date()
// 	y2, m2, d2 := time.Unix(int64(date2), 0).Date()
// 	return y1 == y2 && m1 == m2 && d1 == d2
// }

func (p CoffeeCmd) isReplyOnMyMessage(callback tgbotapi.CallbackQuery) *tgbotapi.Message {
	if callback.Message != nil {
		if msg, ok := p.myMessages[callback.Message.MessageID]; ok {
			return &msg
		}
	}
	return nil
}

func (p *CoffeeCmd) onCallback(callback tgbotapi.CallbackQuery, myMessage tgbotapi.Message) {

	var ID string
	var Button string
	var drink *Drink

	//buttons: OK, Cancel
	split := strings.Split(callback.Data, "::")
	switch len(split) {
	case 2:
		Button = split[0]
		ID = split[1]
		drink = p.entry.getDrinkByID(ID)
	case 1:
		ID = split[0]
		drink = p.entry.getDrinkByID(ID)
	}
	if drink == nil {
		p.notifyUser(callback, "Something went wrong, sorry...")
	} else {
		switch Button {
		case "": //no button, drink has been selected!
			p.onSelectDrink(callback, *drink)

		case btnOK:
			p.onConfirmDrink(callback, *drink)

		case btnCANCEL:
			p.updateMessage(callback, "Request aborted...")
			p.removeInlineKeyboard(callback)
		}
	}
}

func (p *CoffeeCmd) onConfirmDrink(callback tgbotapi.CallbackQuery, drink Drink) {
}

func (p *CoffeeCmd) onSelectDrink(callback tgbotapi.CallbackQuery, drink Drink) {

	log.Printf("Selected: %+v", drink)

	if drink.Entry.Items == nil {
		p.notifyUser(callback, "Good choice!")
		p.removeInlineKeyboard(callback)
		p.updateMessageWithMarkdown(callback, fmt.Sprintf("Please confirm your choice:\n*%s*\nPrice: £%f", drink.ID, drink.Price))
		p.updateInlineKeyboard(callback, confirmDrinkKeyboard(drink))

	} else {
		//next question
		p.updateMessage(callback, drink.Entry.Question)
		p.updateInlineKeyboard(callback, chooseOneDrinkKeyboard(drink.Entry.Items, true))
	}
}

// func buildTree(parent *Drink, drinks []Drink, tree map[string]node) {
// 	for _, item := range drinks {
// 		tree[item.Display] = node{
// 			item:   &item,
// 			parent: parent,
// 		}
// 		buildTree(&item, item.SubItems, tree)
// 	}
// }

func initCoffeeCmd(bot Bot, message tgbotapi.Message) *CoffeeCmd {
	newCmd := &CoffeeCmd{
		Bot:        bot,
		initialMsg: message,
		chatID:     message.Chat.ID,
		myMessages: make(map[int]tgbotapi.Message),
	}

	const filePath = "./data/benugo.json"
	menu, err := loadBewerages(filePath)
	if err != nil {
		dir, _ := os.Getwd()
		log.Printf("%s", dir)
		log.Fatal(err)
	}
	newCmd.entry = menu.Entry
	// newCmd.allDrinks = make(map[string]node, 0)
	// buildTree(nil, newCmd.drinks, newCmd.allDrinks)
	all := newCmd.entry.getAllEntries()
	log.Printf("Bewerages loaded from file: %s, available items: %d", filePath, len(all))
	return newCmd
}
