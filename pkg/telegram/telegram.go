package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"gopkg.in/telebot.v3"

	"github.com/mykysha/StickerifierBot/pkg/photoconverter"
)

var errEmptyToken = errors.New("empty token")

type Client struct {
	bot *telebot.Bot
}

func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("failed to create Telegram client: %w", errEmptyToken)
	}

	settings := telebot.Settings{
		Token:       os.Getenv("TELEGRAM_TOKEN"),
		Synchronous: true,
		Verbose:     true,
	}

	tgBot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, fmt.Errorf("could not create bot %w", err)
	}

	client := &Client{
		bot: tgBot,
	}

	tgBot.Handle(telebot.OnMedia, func(context telebot.Context) error {
		return client.mediaHandler(context.Message())
	})

	tgBot.Handle(telebot.OnPhoto, func(context telebot.Context) error {
		return client.photoHandler(context.Message())
	})

	return client, nil
}

func (c *Client) errorRespond(message *telebot.Message) {
	if _, err := c.bot.Send(message.Chat, "Something went wrong. Please try again later."); err != nil {
		return
	}
}

func (c *Client) mediaHandler(message *telebot.Message) error {
	if _, err := c.bot.Send(message.Chat, "Media message received"); err != nil {
		return fmt.Errorf("could not send message %w", err)
	}

	file, err := c.bot.File(message.Media().MediaFile())
	if err != nil {
		c.errorRespond(message)

		return fmt.Errorf("could not download file %w", err)
	}

	if err = c.fileManipulator(file, message); err != nil {
		c.errorRespond(message)

		return fmt.Errorf("could not convert and resise file %w", err)
	}

	return nil
}

func (c *Client) photoHandler(message *telebot.Message) error {
	if _, err := c.bot.Send(message.Chat, "Photo message received"); err != nil {
		c.errorRespond(message)

		return fmt.Errorf("could not send message %w", err)
	}

	photo := message.Photo

	file, err := c.bot.File(photo.MediaFile())
	if err != nil {
		c.errorRespond(message)

		return fmt.Errorf("could not download file %w", err)
	}

	if err = c.fileManipulator(file, message); err != nil {
		c.errorRespond(message)

		return fmt.Errorf("could not convert and resise file %w", err)
	}

	return nil
}

func (c *Client) fileManipulator(file io.ReadCloser, message *telebot.Message) error {
	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(file); err != nil {
		c.errorRespond(message)

		return fmt.Errorf("could not read file %w", err)
	}

	png, err := photoconverter.ConvertToPNG(buf.Bytes())
	if err != nil {
		c.errorRespond(message)

		return fmt.Errorf("could not convert photo to png %w", err)
	}

	pngPhoto := &telebot.Document{File: telebot.FromReader(bytes.NewReader(png))}

	pngPhoto.FileName = time.Now().Format("2006-01-02_15-04-05") + ".png"
	pngPhoto.MIME = "image/png"
	pngPhoto.Caption = "Here is your photo! Now as PNG!"

	if _, err = c.bot.Send(message.Chat, pngPhoto); err != nil {
		return fmt.Errorf("could not send photo %w", err)
	}

	return nil
}

func (c *Client) HandleWebHook(req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	update := new(telebot.Update)

	if err := json.Unmarshal([]byte(req.Body), update); err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("could not decode incoming update %w", err)
	}

	c.bot.ProcessUpdate(*update)

	return events.APIGatewayProxyResponse{
		Body:       "OK",
		StatusCode: http.StatusOK,
	}, nil
}
