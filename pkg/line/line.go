package line

import (
	"log"
	"store/pkg/config"
	"store/pkg/sheet"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

type bot struct {
	client       *linebot.Client
	sheetService sheet.Service
}

func NewBot(sheetService sheet.Service, cfg *config.AppConfig) *bot {
	client, err := linebot.New(cfg.Line.LineChannelSecret, cfg.Line.LineChannelToken)
	if err != nil {
		log.Fatal(err)
	}

	return &bot{client: client, sheetService: sheetService}
}

func (b *bot) Callback(c *gin.Context) {
	b.handleEvents(b.getEvents(c))
}

func (b *bot) handleEvents(events []*linebot.Event) {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			b.replyMessage(event)
		}
	}
}

func (b *bot) getEvents(c *gin.Context) []*linebot.Event {
	events, err := b.client.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.Writer.WriteHeader(400)
		} else {
			c.Writer.WriteHeader(500)
		}

		return nil
	}

	return events
}
