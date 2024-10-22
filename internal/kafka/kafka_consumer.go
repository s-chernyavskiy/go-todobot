package kafka

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"go-todobot/internal/domain/events"
	"go-todobot/internal/storage"
	"go-todobot/lib/e"
	"log"
	"strings"
)

type Consumer struct {
	r *kafka.Reader
}

func (c *Consumer) Process(event events.Event) error {
	log.Println("processing event:", event)
	switch event.Type {
	case events.Message:
		return c.processMessage(event)
	default:
		return ErrUnknownEventType
	}
}

func (c *Consumer) processMessage(event events.Event) error {
	meta, err := meta(event)
	// {1 xxsd map[ChatID:430746829 Username:jugglyyfe]}
	// {1 asd {430746829 jugglyyfe}}

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
	// case StartCmd:
	// 	return c.sendHello(chatID)
	// case ListCmd:
	// 	return c.listTasks(chatID, username)
	// case AddCmd:
	// 	return c.addTask(text, chatID, username)
	// case RemoveCmd:
	// 	return c.removeTask(text, chatID, username)
	// default:
	// 	return c.tg.SendMessage(chatID, msgUnknownCommand)
	}
	return nil
}

func (c *Consumer) Close() error {
	return c.r.Close()
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		r: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: groupID,
			Topic:   topic,
		}),
	}
}

const (
	ListCmd   = "/list"
	HelpCmd   = "/help"
	StartCmd  = "/start"
	AddCmd    = "/add"
	RemoveCmd = "/remove"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
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
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	case ListCmd:
		return p.listTasks(chatID, username)
	case AddCmd:
		return p.addTask(text, chatID, username)
	case RemoveCmd:
		return p.removeTask(text, chatID, username)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) removeTask(text string, chatID int, username string) (err error) {
	defer func() { err = e.NilWrap("cant add task", err) }()
	log.Println("start removing task:", text)

	t := &storage.Task{
		UserID: username,
		Name:   text,
	}

	log.Println("checking if task exists:", text)
	isExists, err := p.storage.IfExistsTask(t)
	if err != nil {
		log.Println("task exist check error:", text)
		return err
	}

	if !isExists {
		log.Println("task doesnt exist:", text)
		return p.tg.SendMessage(chatID, msgTaskDoesntExist)
	}

	errCh := make(chan error)

	go func() {
		errCh <- p.storage.RemoveTask(t)
	}()

	if err := <-errCh; err != nil {
		log.Println("task removing error:", text)
		return err
	}

	log.Println("sending message:", text)
	if err := p.tg.SendMessage(chatID, msgDeleted); err != nil {
		log.Println("sending message error:", text)
		return err
	}

	return nil
}

func (p *Processor) addTask(text string, chatID int, username string) (err error) {
	defer func() { err = e.NilWrap("cant add task", err) }()
	log.Println("adding task:", text)
	t := &storage.Task{
		UserID: username,
		Name:   text,
	}

	isExists, err := p.storage.IfExistsTask(t)
	if err != nil {
		log.Println("task exist check error:", text)
		return err
	}
	if isExists {
		log.Println("task exists:", text)
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	errCh := make(chan error)

	go func() {
		errCh <- p.storage.AddTask(t)
	}()

	if err := <-errCh; err != nil {
		log.Println("task adding error:", text)
		return err
	}

	if err := p.tg.SendMessage(chatID, msgCreated); err != nil {
		log.Println("sending message error:", text)
		return err
	}

	return nil
}

func (p *Processor) listTasks(chatID int, username string) (err error) {
	defer func() { err = e.NilWrap("cant list tasks", err) }()

	tasks, err := p.storage.ListTasks(username)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		err := p.tg.SendMessage(chatID, msgNoTaskFound)
		if err != nil {
			return err
		}

		return nil
	}

	for i, task := range tasks {
		err := p.tg.SendMessage(chatID, fmt.Sprintf("%d: %s", i+1, task))
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}
