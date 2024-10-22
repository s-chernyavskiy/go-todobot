package telegram

import (
	"errors"
	"go-todobot/internal/domain/events"
	"go-todobot/internal/domain/models"
	"go-todobot/internal/storage"
	"go-todobot/internal/telegram"
	"go-todobot/lib/e"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(tg *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      tg,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
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

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return ErrUnknownEventType
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	// {1 xxsd map[ChatID:430746829 Username:jugglyyfe]}
	// {1 asd {430746829 jugglyyfe}}

	if err != nil {
		return e.Wrap("cant process message", err)
	}
	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("cant process message", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)

	if !ok {
		return Meta{}, e.Wrap("cant get meta", ErrUnknownMetaType)
	}

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

func fetchType(upd models.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func fetchText(upd models.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}
