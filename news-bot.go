package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//NewsBot :
type NewsBot struct {
	Bot
	MessageHandler

	Feeds map[string]string
}

//RSS feed structure
type RSS struct {
	Items []Item `xml:"channel>item"`
}

//Item : element of RSS feed
type Item struct {
	URL   string `xml:"guid"`
	Title string `xml:"title"`
}

func (bot NewsBot) getFeed(tag string) string {
	key := strings.ToLower(tag)
	for k, v := range bot.Feeds {
		if key == strings.ToLower(k) {
			return v
		}
	}
	return ""
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

//ProcessUpdate : handling messages for the bot
func (bot NewsBot) ProcessUpdate(update tgbotapi.Update) {

	if url := bot.getFeed(update.Message.Text); url != "" {
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

func initNewsBot(token string, debug bool) *NewsBot {
	if token == "" {
		return nil
	}

	bot := &NewsBot{
		Bot: Bot{
			"News Bot",
			initBotAPI(token, debug),
		},
		Feeds: map[string]string{
			"Habr":     "https://habrahabr.ru/rss/best/",
			"Learn Go": "https://blog.learngoprogramming.com/feed",
			"BBC News": "http://newsrss.bbc.co.uk/rss/newsonline_uk_edition/front_page/rss.xml",
		},
	}

	bot.logBotDetails()
	return bot
}
