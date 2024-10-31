package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/mykysha/StickerifierBot/pkg/telegram"
)

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")

	if token == "" {
		panic("TELEGRAM_TOKEN is not set")
	}

	client, err := telegram.NewClient(token)
	if err != nil {
		panic(err)
	}

	lambda.Start(client.HandleWebHook)
}
