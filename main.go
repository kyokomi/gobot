package main

import (
	"fmt"
	"os"

	"encoding/json"
	"log"

	"flag"

	docomo "github.com/kyokomi/go-docomo"
	"github.com/kyokomi/nepu-bot/victor"
)

var logger = log.New(os.Stderr, "nepu-bot", log.Llongfile)

func main() {

	var apiKey string
	flag.StringVar(&apiKey, "APIKEY", "", "docomo developerで登録したAPIKEYをして下さい")
	flag.Parse()

	if apiKey == "" {
		//log.Fatalln("APIKEYを指定して下さい")
		apiKey = os.Getenv("DOCOMO_APIKEY")
	}

	bot := victor.New(victor.Config{
		Name:         "いーすん",
		ChatAdapter:  "slack",
		StoreAdapter: "memory",
		HTTPAddr:     ":80",
	})

	d := docomo.New(apiKey)

	bot.HandleCommandFunc("hello|hi|howdy", (victor.HandlerFunc)(func(s victor.State) {
		s.Chat().Send(s.Message().ChannelID(), fmt.Sprintf("Hello, %s", s.Message().UserName()))
	}))

	bot.HandleCommandFunc(".*", (victor.HandlerFunc)(func(s victor.State) {
		res, err := d.SendZatsudan(s.Message().UserName(), s.Message().Text())
		if err != nil {
			logger.Println(err)
			return
		}

		var resMap map[string]string
		if err := json.Unmarshal(res, &resMap); err != nil {
			logger.Println("Unmarshal ", err)
			return
		}

		// Send Slack
		s.Chat().Send(s.Message().ChannelID(), resMap["utt"])
	}))

	bot.Run()
}
