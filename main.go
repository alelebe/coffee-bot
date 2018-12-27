package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Vars :
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
	var err error

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	var newsCh tgbotapi.UpdatesChannel

	if p.news != nil {
		p.news.RemoveWebhook()
		newsCh, err = p.news.GetUpdatesChan(u)
		if err != nil {
			panic(err)
		}
	}

	if newsCh == nil {
		log.Println("Nothing to do :-( -> check logs")
		return
	}
	// читаем обновления из канала
	for {
		select {
		case update := <-newsCh:
			MessageHandler(p.news).ProcessUpdate(update)
		}
	}
}

func (p Program) configureHook(bot Bot, router *gin.Engine, handler MessageHandler) bool {
	var err error

	err = bot.setupWebhook(p.baseURL, router, handler)
	if err != nil {
		log.Printf("Fail to set WebHook for '%s': %v\n", bot.Name, err)
		return false
	}
	return true
}

func (p Program) runRouter() {
	var err error
	// run webHooks on Gin router
	router := gin.New()
	router.Use(gin.Logger())

	configured := false
	if p.news != nil {
		configured = configured ||
			p.configureHook(p.news.Bot, router, MessageHandler(p.news))
	}

	if configured {
		err = router.Run(":" + p.port)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Nothing to do :-( -> check logs")
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
	govendor fetch github.com/go-telegram-bot-api/telegram-bot-api

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
	heroku local -e .env.test
	hrekou ci:debug

HOWTO: stop the app
	heroku ps:scale web=0

HOWTO: start the app
	$ git push heroku master
	$ heroku ps:scale web=1
	$ heroku open

HOWTO: change app variable
	$ heroku config:set GIN_MODE=release

HOWTO: debug
https://github.com/Microsoft/vscode-go/wiki/Debugging-Go-code-using-VS-Code
	>> Install Delve
	$ xcode-select --install
	$ go get -u github.com/derekparker/delve/cmd/dlv
*/
