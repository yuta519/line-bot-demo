package main

import (
	"fmt"

	"github.com/yuta519/line-bot-demo/infra"
)

func main() {
	token := infra.FetchChannelAccessToken()
	fmt.Println(token)
}
