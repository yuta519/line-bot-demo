package main

import (
	"log"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/yuta519/line-bot-demo/infra"
)

func main() {
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		infra.FetchChannelAccessToken(),
	)
	if err != nil {
		log.Fatal(err)
	}

	message := linebot.NewTextMessage("hogehogehoge")
	if _, err := bot.BroadcastMessage(message).Do(); err != nil {
		log.Fatal(err)
	}
}
