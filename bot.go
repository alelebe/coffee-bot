package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//MessageHandler : processing messages from BotAPI
type MessageHandler interface {
	ProcessUpdate(tgbotapi.Update)
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

func (bot Bot) dispatchMessage(update tgbotapi.Update, handler MessageHandler) {
	if update.Message == nil { // ignore any non-Message Updates
		log.Printf("UPDATE: %+v\n", update)
		return
	}

	// log.Printf("Date: %v\n", time.Unix(int64(update.Message.Date), 0))
	log.Printf("From %+v (%s): %+v\n", update.Message.From, update.Message.Chat.Type, update.Message.Text)

	handler.ProcessUpdate(update)
}

func (bot Bot) ginWebhook(c *gin.Context, handler MessageHandler) {
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

func (bot Bot) setupWebhook(baseURL string, router *gin.Engine, handler MessageHandler) error {
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
