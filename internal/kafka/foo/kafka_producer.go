package foo

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"go-todobot/internal/domain/events"
	"go-todobot/internal/domain/models"
	"go-todobot/internal/telegram"
	"go-todobot/lib/e"
)

type Producer struct {
	writer *kafka.Writer
	tg     *telegram.Client
	offset int
}

type Meta struct {
	ChatID   int
	Username string
}

func (p *Producer) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return errors.New("unknown event")
	}
}

func (p *Producer) processMessage(event events.Event) error {
	meta, err := meta(event)
	_ = meta
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Value: []byte("hello world"),
	}

	err = p.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		return err
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)

	if !ok {
		return Meta{}, e.Wrap("cant get meta", errors.New("unknown meta"))
	}

	return res, nil
}

func New(tg *telegram.Client, kafkaURL, topic string) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{
		tg:     tg,
		writer: w,
	}
}

func (p *Producer) getUpdates(limit int) ([]events.Event, error) {
	upd, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("cant get events", err)
	}

	if len(upd) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(upd))

	for _, u := range upd {
		res = append(res, event(u))
	}

	p.offset = upd[len(upd)-1].ID + 1

	return res, nil
}

func event(upd models.Update) events.Event {
	t := fetchType(upd)

	res := events.Event{
		Type: t,
		Text: fetchText(upd),
	}

	if t == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd models.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd models.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
