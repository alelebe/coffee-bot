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

	log.Printf("Authorized on account %s", botAPI.Self.UserName)
	return botAPI
}

func dispatchMessage(c *gin.Context, handler MessageHandler) {
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

	if update.Message == nil { // ignore any non-Message Updates
		return
	}
	handler.ProcessUpdate(update)
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
		dispatchMessage(c, handler)
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

	if bot.Debug {
		log.Println("Bot debug mode is ON")
		log.Printf("Bot buffer: %d\n", bot.Buffer)
	}
}
