package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//UpdateHandler : processing messages from BotAPI
type UpdateHandler interface {
	ProcessMessage(tgbotapi.Message)
	ProcessCallback(tgbotapi.CallbackQuery)
}

//Bot :
type Bot struct {
	Name string
	*tgbotapi.BotAPI
}

func initBotAPI(token string, debug bool) *tgbotapi.BotAPI {
	var err error

	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	botAPI.Debug = debug
	return botAPI
}

func (bot Bot) logBotDetails() {
	log.Printf("%s: Authorized on account %s\n", bot.Name, bot.Self.UserName)

	if bot.Debug {
		log.Printf("%s: Debug mode is enabled\n", bot.Name)
		log.Printf("%s: Buffer: %d\n", bot.Name, bot.Buffer)
	}
}

func (bot Bot) dispatchMessage(update tgbotapi.Update, handler UpdateHandler) {
	if update.Message != nil {
		// log.Printf("Date: %v\n", time.Unix(int64(update.Message.Date), 0))
		log.Printf("From %+v (%s): %s\n", update.Message.From, update.Message.Chat.Type, update.Message.Text)

		handler.ProcessMessage(*update.Message)

	} else if update.CallbackQuery != nil {

		log.Printf("Callback >> From %+v: %s\n", update.CallbackQuery.From, update.CallbackQuery.Data)
		handler.ProcessCallback(*update.CallbackQuery)

	} else {
		// ignore any non-Message Updates
		log.Printf("UPDATE: %+v\n", update)
		return
	}

}

func (bot Bot) ginWebhook(c *gin.Context, handler UpdateHandler) {
	defer c.Request.Body.Close()

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var update tgbotapi.Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Println(err)
		return
	}

	bot.dispatchMessage(update, handler)
}

func (bot Bot) setupWebhook(baseURL string, router *gin.Engine, handler UpdateHandler) error {
	var err error

	base, err := url.Parse(baseURL)
	if err != nil {
		return err
	}
	url, err := url.Parse("/bot" + bot.Token)
	if err != nil {
		return err
	}
	hookURL := base.ResolveReference(url)

	// this perhaps should be conditional on GetWebhookInfo()
	// only set webhook if it is not set properly
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(hookURL.String()))
	if err != nil {
		return err
	}

	router.POST("/bot"+bot.Token, func(c *gin.Context) {
		bot.ginWebhook(c, handler)
	})
	return nil
}

func (bot Bot) logWebhookDetails() {
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Web hook: %s\n", info.URL)

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s\n", info.LastErrorMessage)
	}

	if info.PendingUpdateCount != 0 {
		log.Printf("Pending updates: %d\n", info.PendingUpdateCount)
	}
}

func (bot Bot) notifyUser(callback tgbotapi.CallbackQuery, text string) {
	bot.AnswerCallbackQuery(
		tgbotapi.NewCallback(
			callback.ID,
			text,
		))
}

func (bot Bot) alertUser(callback tgbotapi.CallbackQuery, text string) {
	bot.AnswerCallbackQuery(
		tgbotapi.NewCallbackWithAlert(
			callback.ID,
			text,
		))
}

func (bot Bot) sendToChat(chatID int64, text string) (tgbotapi.Message, error) {
	return bot.Send(
		tgbotapi.NewMessage(
			chatID,
			text,
		))
}

func (bot Bot) replyToMessage(message tgbotapi.Message, text string) (tgbotapi.Message, error) {
	return bot.Send(
		tgbotapi.NewMessage(
			message.Chat.ID,
			text,
		))
}

func (bot Bot) replyToMessageWithInlineKeyboard(message tgbotapi.Message, text string,
	keyboard [][]tgbotapi.InlineKeyboardButton) (tgbotapi.Message, error) {

	msg := tgbotapi.NewMessage(
		message.Chat.ID,
		text,
	)
	msg.ParseMode = "Markdown"
	msg.BaseChat.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)

	return bot.Send(msg)
}

func (bot Bot) updateMessage(callback tgbotapi.CallbackQuery, text string) {
	msg := tgbotapi.NewEditMessageText(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		text,
	)
	msg.ParseMode = "Markdown"

	bot.Send(msg)
}

func (bot Bot) updateInlineKeyboard(callback tgbotapi.CallbackQuery, keyboard [][]tgbotapi.InlineKeyboardButton) {
	bot.Send(
		tgbotapi.NewEditMessageReplyMarkup(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			tgbotapi.NewInlineKeyboardMarkup(
				keyboard...,
			),
		))
}

func (bot Bot) removeInlineKeyboard(callback tgbotapi.CallbackQuery) {
	bot.Send(
		tgbotapi.NewEditMessageReplyMarkup(
			callback.Message.Chat.ID, callback.Message.MessageID,
			tgbotapi.NewInlineKeyboardMarkup(make([]tgbotapi.InlineKeyboardButton, 0)),
		))
}
