package clarifaitelegram

import (
	"errors"
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/clarifai/clarifai-go"
	"log"
	"os"
	"strings"
)

var (
	errClarifaiIDEmpty     = errors.New("Clarifai ID is empty string")
	errClarifaiSecretEmpty = errors.New("Clarifai Secret is empty string")
	errTelegramTokenEmpty  = errors.New("Telegram token is empty string")
)

// Client represents ClarifaiID, ClarifaiSecret and TelegramToken
type Client struct {
	ClarifaiID     string
	ClarifaiSecret string
	TelegramToken  string
}

func (cl *Client) getTags(url string) ([]string, error) {
	if cl.ClarifaiID == "" {
		return nil, errClarifaiIDEmpty
	}
	if cl.ClarifaiSecret == "" {
		return nil, errClarifaiSecretEmpty
	}

	client := clarifai.NewClient(cl.ClarifaiID, cl.ClarifaiSecret)
	_, err := client.Info()
	if err != nil {
		return nil, err
	}

	urls := []string{url}
	tagdata, err := client.Tag(urls, nil)

	if err != nil {
		return nil, err
	}

	return tagdata.Results[0].Result.Tag.Classes, nil
}

// LoadFromEnv provides loading secrets and tokens from environment variables
// TELEGRAM_TOKEN
// CLARIFAI_ID
// CLARIFAI_SECRET
func (cl *Client) LoadFromEnv() {
	cl.TelegramToken = os.Getenv("TELEGRAM_TOKEN")
	cl.ClarifaiID = os.Getenv("CLARIFAI_ID")
	cl.ClarifaiSecret = os.Getenv("CLARIFAI_SECRET")
}

// Start provides entry point for bot
func (cl *Client) Start() {
	if cl.TelegramToken == "" {
		log.Fatal(errTelegramTokenEmpty)
	}

	bot, err := tgbotapi.NewBotAPI(cl.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	err = bot.UpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for update := range bot.Updates {

			text := update.Message.Text
			if text == "/link" {
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Paste your link to the image")
				bot.SendMessage(msg)
			}

			if strings.HasPrefix(text, "http") {
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID
				result, err := cl.getTags(text)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Can't receive data from Clarifai")
					msg.ReplyToMessageID = update.Message.MessageID
					bot.SendMessage(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%+v\n", result))
					msg.ReplyToMessageID = update.Message.MessageID
					bot.SendMessage(msg)
				}

			}
			fmt.Println(update)
		}
	}
}
