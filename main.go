package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	bot *tgbotapi.BotAPI
)

var rss = map[string]string{
	"Habr": "https://habrahabr.ru/rss/best/",
}

type RSS struct {
	Items []Item `xml:"channel>item"`
}

type Item struct {
	URL   string `xml:"guid"`
	Title string `xml:"title"`
}

func getNews(url string) (*RSS, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	rss := new(RSS)
	err = xml.Unmarshal(body, rss)
	if err != nil {
		return nil, err
	}

	return rss, nil
}

func initTelegramBot() {
	var err error

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Println("$TOKEN must be set")
	}

	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	debugBot := os.Getenv("DEBUG_BOT")
	if debugBot != "" && debugBot != "0" {
		bot.Debug = true
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func removeWebhook() {
	bot.RemoveWebhook()
}

func setupWebhook() {
	var err error

	base := os.Getenv("BASE_URL")
	if base == "" {
		log.Fatal("$BASE_URL must be set")
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		log.Fatal(err)
	}
	url, err := url.Parse("/bot" + bot.Token)
	if err != nil {
		log.Fatal(err)
	}
	hookURL := baseURL.ResolveReference(url)

	// this perhaps should be conditional on GetWebhookInfo()
	// only set webhook if it is not set properly
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(hookURL.String()))
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	log.Printf("Pending updates: %d\n", info.PendingUpdateCount)
}

func webhookHandler(c *gin.Context) {
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
	processUpdate(update)
}

func processUpdate(update tgbotapi.Update) {

	if update.Message == nil { // ignore any non-Message Updates
		return
	}
	// to monitor changes run: heroku logs --tail
	log.Printf("From %+v: %+v\n", update.Message.From, update.Message.Text)

	if url, ok := rss[update.Message.Text]; ok {
		rss, err := getNews(url)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"sorry, error happend",
			))
		}
		for _, item := range rss.Items {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				item.URL+"\n"+item.Title,
			))
		}
	} else {
		bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			`there is only Habr feed availible`,
		))
	}
}

func runLocally() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	// читаем обновления из канала
	for {
		select {
		case update := <-updates:
			processUpdate(update)
		}
	}
}

func runRouter() {
	var err error

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// gin router
	router := gin.New()
	router.Use(gin.Logger())

	router.POST("/bot"+bot.Token, webhookHandler)

	err = router.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}

func main() {

	// telegram
	initTelegramBot()

	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}
	log.Println("Bot server started in the " + env + " mode")

	if strings.ToLower(env) == "local" {
		removeWebhook()
		runLocally()
	} else {
		setupWebhook()
		runRouter()
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
