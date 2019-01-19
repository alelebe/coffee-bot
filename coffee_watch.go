package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//CoffeeWatch :
type CoffeeWatch struct {
	Bot
	initialMsg tgbotapi.Message
	chatID     int64

	myMessages map[int]tgbotapi.Message
}

func initCoffeeWatch(bot Bot, message tgbotapi.Message) *CoffeeWatch {

	newChat := &CoffeeWatch{
		Bot:        bot,
		initialMsg: message,
		chatID:     message.Chat.ID,

		myMessages: make(map[int]tgbotapi.Message),
	}
	return newChat
}

func whatWouldYouLikeToDo(amIwatcher bool, numOfWatchers int) string {

	var result string
	if amIwatcher {
		result = "Allright. You are already a coffee watcher...\n"
		numOfWatchers--
	}

	if numOfWatchers > 1 {
		s := ""
		if numOfWatchers > 2 {
			s = "s"
		}

		if amIwatcher {
			result += fmt.Sprintf("There are also %d other%s watching coffee requests\n", numOfWatchers, s)
		} else {
			result += fmt.Sprintf("There are %d other%s watching coffee requests\n", numOfWatchers, s)
		}
	}

	return result + "Please tell me what would you like to do?"
}

func startOrStop(amIwatcher bool) [][]tgbotapi.InlineKeyboardButton {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	if amIwatcher {
		keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(btnSTOP, btnSTOP),
			tgbotapi.NewInlineKeyboardButtonData(btnCANCEL, btnCANCEL),
		))
	} else {
		keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(btnWATCH, btnWATCH),
			tgbotapi.NewInlineKeyboardButtonData(btnCANCEL, btnCANCEL),
		))
	}
	return keyboard
}

func (p *CoffeeWatch) start() {

	amIwatcher, watchers := amIcoffeeWatcher(p.initialMsg.From.ID)

	sent, err := p.replyToMessageWithInlineKeyboard(
		p.initialMsg, whatWouldYouLikeToDo(amIwatcher, len(watchers)),
		startOrStop(amIwatcher),
	)
	if err == nil {
		p.myMessages[sent.MessageID] = sent
	}
	log.Printf("Coffee Watch: new request with msgId: %d", sent.MessageID)
}

func (p CoffeeWatch) isReplyOnMyMessage(callback tgbotapi.CallbackQuery) bool {
	if callback.Message != nil {
		if _, ok := p.myMessages[callback.Message.MessageID]; ok {
			return true
		}
	}
	return false
}

func (p *CoffeeWatch) onCallback(callback tgbotapi.CallbackQuery) {

	button := callback.Data
	switch button {
	default:
		return

	case btnWATCH:
		p.startWatching(callback)

	case btnSTOP:
		p.stopWatching(callback)

	case btnCANCEL:
		p.removeInlineKeyboard(callback)
		p.updateMessage(callback, "Allright. You'll tell me what to do later")
	}

	//remove message from my queue
	delete(p.myMessages, callback.Message.MessageID)
}

func (p *CoffeeWatch) startWatching(callback tgbotapi.CallbackQuery) {
	p.removeInlineKeyboard(callback)

	err := addCoffeeWatcher(CoffeeWatcher{
		UserID:   p.initialMsg.From.ID,
		UserName: p.initialMsg.From.FirstName,
		ChatID:   p.initialMsg.Chat.ID,
	})
	if err == nil {
		p.updateMessage(callback, "I've added you to the list of watchers.\nEnjoy tracking coffee chat...")
	} else {
		p.updateMessage(callback, "I'm sorry.. Something went wrong")
	}
}

func (p *CoffeeWatch) stopWatching(callback tgbotapi.CallbackQuery) {
	p.removeInlineKeyboard(callback)

	err := removeCoffeeWatcher(CoffeeWatcher{
		UserID:   p.initialMsg.From.ID,
		UserName: p.initialMsg.From.FirstName,
		ChatID:   p.initialMsg.Chat.ID,
	})

	if err == nil {
		p.updateMessage(callback, "I've removed you from the list of watchers.")
	} else {
		p.updateMessage(callback, "I'm sorry.. Something went wrong")
	}
}
