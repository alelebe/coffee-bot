package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Vars struct {
	port    string
	baseURL string

	debugBot bool
	mode     string
}

//Program :
type Program struct {
	Vars
	news *NewsBot
}

func (p Program) isLocal() bool {
	return p.mode == "LOCAL" || p.mode == ""
}

func (p Program) runLongPooling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	p.news.RemoveWebhook()
	newsCh, err := p.news.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	// читаем обновления из канала
	for {
		select {
		case update := <-newsCh:
			MessageHandler(p.news).ProcessUpdate(update)
		}
	}
}

func (p Program) runRouter() {
	// run webHooks on Gin router
	router := gin.New()
	router.Use(gin.Logger())

	p.news.setupWebhook(p.baseURL, router, MessageHandler(p.news))

	err := router.Run(":" + p.port)
	if err != nil {
		log.Fatal(err)
	}
}

func getVar(env string) string {
	value := os.Getenv(env)
	if value == "" {
		log.Printf("$%s must be set\n", env)
	}
	return value
}
func getOptVar(env string, defValue string) string {
	value := strings.ToLower(os.Getenv("ENV"))
	if value == "" {
		value = defValue
	}
	return value
}

func initVars() Vars {
	port := getVar("PORT")
	baseURL := getVar("BASE_URL")

	debugStr := getOptVar("DEBUG_BOT", "0")

	debug := false
	if debugStr != "" && debugStr != "0" {
		debug = true
	}

	mode := strings.ToUpper(getVar("ENV"))

	return Vars{
		port:     port,
		baseURL:  baseURL,
		debugBot: debug,
		mode:     mode,
	}
}

func main() {
	p := Program{
		Vars: initVars(),
	}
	log.Println("Started in " + p.mode + " mode")

	// construct Telergam Bots
	p.news = initNewsBot(getVar("NEWS_TOKEN"), p.debugBot)

	if p.isLocal() {
		p.runLongPooling()
	} else {
		p.runRouter()
	}
}

/*
	go get gopkg.in/telegram-bot-api.v4
	heroku git:remote -a fx-coffee-bot

	govendor init
	govendor fetch github.com/gin-gonic/gin

	heroku plugins:install @heroku-cli/plugin-manifest
	heroku manifest:create

	git push heroku master
	heroku logs --tail

	Telegram - speak to 'BotFather'
		name:		alelebeGoHabr
		username:	alelebe_habr_bot

		Use this token to access the HTTP API:
		***REMOVED***

	https://dashboard.ngrok.com/get-started
	./ngrok http 8080
		==> update WebhookURL

	go install
	heroku local

HOWTO: stop the app
	heroku ps:scale web=0

	$ git push heroku master
	$ heroku ps:scale web=1
	$ heroku open
*/
