package main

import (
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/yuta519/line-bot-demo/infra"
)

func main() {
	access_token := infra.FetchChannelAccessToken()
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		access_token,
	)
	if err != nil {
		log.Fatal(err)
	}

	message := linebot.NewTextMessage("MessagingAPIテスト")
	if _, err := bot.BroadcastMessage(message).Do(); err != nil {
		log.Fatal(err)
	}

	infra.RevokeAccessToken(access_token)

	http.HandleFunc("/", testHandler(bot, "aaaa"))
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func testHandler(bot *linebot.Client, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi, from Service: " + name))
	}
}
