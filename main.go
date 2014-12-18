package main

import (
	"os"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	docomo "github.com/kyokomi/go-docomo"
	"github.com/zenazn/goji"
)

type Config struct {
	Name,
	ChatAdapter,
	StoreAdapter,
	HTTPAddr string
}

type Message struct {
	userID, userName, channelID, channelName, text string
}

type OutgoingMessage struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

var logger = log.New(os.Stderr, "nepu-bot", log.Llongfile)

func main() {

	bot := Config{
		Name:         "いーすん",
		ChatAdapter:  "slack",
		StoreAdapter: "memory",
		HTTPAddr:     os.Getenv("PORT"),
	}
	team := os.Getenv("VICTOR_SLACK_TEAM")
	token := os.Getenv("VICTOR_SLACK_TOKEN")
	d := docomo.New(os.Getenv("DOCOMO_APIKEY"))

	goji.Post("/hubot/slack-webhook", func(_ http.ResponseWriter, r *http.Request) {
		m := Message{
			userID:      r.PostFormValue("user_id"),
			userName:    r.PostFormValue("user_name"),
			channelID:   r.PostFormValue("channel_id"),
			channelName: r.PostFormValue("channel_name"),
			text:        r.PostFormValue("text"),
		}

		if !strings.Contains(m.text, bot.Name) {
			return
		}

		res, err := d.SendZatsudan(m.userName, strings.Replace(m.text, bot.Name, "", 1))
		if err != nil {
			logger.Println(err)
			return
		}

		var resMap map[string]string
		if err := json.Unmarshal(res, &resMap); err != nil {
			logger.Println("Unmarshal ", err)
			return
		}

		go Send(bot.Name, team, token, m.channelID, resMap["utt"])
	})
	goji.Serve()
}

func Send(botName, team, token, channelID, msg string) {
	body, err := json.Marshal(&OutgoingMessage{
		Channel:  channelID,
		Username: botName,
		Text:     msg,
	})

	if err != nil {
		log.Println("error sending to chat:", err)
	}

	endpoint := fmt.Sprintf("https://%s.slack.com/services/hooks/hubot?token=%s", team, token)
	http.PostForm(endpoint, url.Values{"payload": {string(body)}})
}
