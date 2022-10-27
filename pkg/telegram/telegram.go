package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nndergunov/stickerify_tgbot/pkg/logger"
	"github.com/nndergunov/stickerify_tgbot/pkg/photoconverter"
	"gopkg.in/telebot.v3"
)

type Client struct {
	log logger.Logger
	bot *telebot.Bot
}

// NewClient creates a new Telegram client.
func NewClient(log logger.Logger, token string) (*Client, error) {
	if log == nil {
		log = logger.NewDefault(os.Stdout, "telegram client")
	}

	if token == "" {
		return nil, fmt.Errorf("failed to create Telegram client: %w", ErrEmptyToken)
	}

	settings := telebot.Settings{
		URL:         "",
		Token:       os.Getenv("TELEGRAM_TOKEN"),
		Updates:     0,
		Poller:      nil,
		Synchronous: true,
		Verbose:     true,
		ParseMode:   "",
		OnError:     nil,
		Client:      nil,
		Offline:     false,
	}

	tgBot, err := telebot.NewBot(settings)
	if err != nil {
		panic(err)
	}

	c := &Client{
		log: log,
		bot: tgBot,
	}

	tgBot.Handle(telebot.OnMedia, func(context telebot.Context) error {
		return c.mediaHandler(context.Message())
	})

	tgBot.Handle(telebot.OnPhoto, func(context telebot.Context) error {
		return c.photoHandler(context.Message())
	})

	return c, nil
}

// mediaHandler handles incoming media messages.
func (c *Client) mediaHandler(m *telebot.Message) error {
	_, err := c.bot.Send(m.Chat, "Media message received")
	if err != nil {
		return fmt.Errorf("could not send message %w", err)
	}

	// extract media from message
	file, err := c.bot.File(m.Media().MediaFile())
	if err != nil {
		c.errorRespond(m)

		return fmt.Errorf("could not download file %w", err)
	}

	err = c.fileManipulator(file, m)
	if err != nil {
		c.errorRespond(m)

		return fmt.Errorf("could not convert and resise file %w", err)
	}

	return nil
}

// photoHandler handles incoming photo messages.
func (c *Client) photoHandler(m *telebot.Message) error {
	_, err := c.bot.Send(m.Chat, "Photo message received")
	if err != nil {
		c.errorRespond(m)

		return fmt.Errorf("could not send message %w", err)
	}

	// extract photo from message
	photo := m.Photo

	// download photo
	file, err := c.bot.File(photo.MediaFile())
	if err != nil {
		c.errorRespond(m)

		return fmt.Errorf("could not download file %w", err)
	}

	err = c.fileManipulator(file, m)
	if err != nil {
		c.errorRespond(m)

		return fmt.Errorf("could not convert and resise file %w", err)
	}

	return nil
}

// fileManipulator further converts and returns photo from photo or media message.
func (c *Client) fileManipulator(file io.ReadCloser, m *telebot.Message) error {
	buf := new(bytes.Buffer)

	_, err := buf.ReadFrom(file)
	if err != nil {
		c.errorRespond(m)

		return fmt.Errorf("could not read file %w", err)
	}

	png, err := photoconverter.ConvertToPNG(buf.Bytes())
	if err != nil {
		c.errorRespond(m)

		return fmt.Errorf("could not convert photo to png %w", err)
	}

	// send photo
	pngPhoto := &telebot.Document{File: telebot.FromReader(bytes.NewReader(png))}

	pngPhoto.FileName = time.Now().Format("2006-01-02_15-04-05") + ".png"
	pngPhoto.MIME = "image/png"
	pngPhoto.Caption = "Here is your photo! Now as PNG!"

	_, err = c.bot.Send(m.Chat, pngPhoto)
	if err != nil {
		return fmt.Errorf("could not send photo %w", err)
	}

	return nil
}

// errorRespond sends error message to user.
func (c *Client) errorRespond(m *telebot.Message) {
	_, err := c.bot.Send(m.Chat, "Something went wrong. Please try again later.")
	if err != nil {
		return
	}
}

// HandleWebHook handles serverless webhook communication.
func (c *Client) HandleWebHook(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	update := new(telebot.Update)

	err := json.Unmarshal([]byte(req.Body), update)
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("could not decode incoming update %w", err)
	}

	c.bot.ProcessUpdate(*update)

	return events.APIGatewayProxyResponse{
		Body:       "OK",
		StatusCode: 200,
	}, nil
}
