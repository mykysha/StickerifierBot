package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nndergunov/stickerify_tgbot/pkg/logger"
	"github.com/nndergunov/stickerify_tgbot/pkg/telegram"
)

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")

	if token == "" {
		panic("TELEGRAM_TOKEN is not set")
	}

	tgLogger := logger.NewDefault(os.Stdout, "telegram")

	client, err := telegram.NewClient(tgLogger, token)
	if err != nil {
		panic(err)
	}

	lambda.Start(client.HandleWebHook)
}
