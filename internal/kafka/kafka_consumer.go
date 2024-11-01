package kafka

import (
	"github.com/segmentio/kafka-go"
	"go-todobot/internal/domain/events"
	"go-todobot/internal/telegram"
	"go-todobot/lib/e"
	"log"
	"strings"
)

type Consumer struct {
	tg *telegram.Client
	r  *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string, tg *telegram.Client) *Consumer {
	return &Consumer{
		tg: tg,
		r: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: groupID,
			Topic:   topic,
		}),
	}
}

func (c *Consumer) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return c.processMessage(event)
	default:
		return ErrUnknownEventType
	}
}

func (c *Consumer) processMessage(event events.Event) error {
	meta, err := meta(event)

	if err != nil {
		return e.Wrap("cant process message", err)
	}
	if err := c.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("cant process message", err)
	}
	return nil
}

func (c *Consumer) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command %s from %s", text, username)

	firstSpace := strings.Index(text, " ")

	var req string
	if firstSpace != -1 {
		req = text[:firstSpace]
		text = text[firstSpace+1:]
	} else {
		req = text
	}

	switch req {
	// case HelpCmd:
	// 	return c.sendHelp(chatID)
	case StartCmd:
		return c.sendHello(chatID)
	// case ListCmd:
	// 	return c.listTasks(chatID, username)
	// case AddCmd:
	// 	return c.addTask(text, chatID, username)
	// case RemoveCmd:
	// 	return c.removeTask(text, chatID, username)
	default:
		return c.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (c *Consumer) Close() error {
	return c.r.Close()
}

const (
	ListCmd   = "/list"
	HelpCmd   = "/help"
	StartCmd  = "/start"
	AddCmd    = "/add"
	RemoveCmd = "/remove"
)
