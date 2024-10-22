package foo

import (
	"context"
	"go-todobot/internal/client/telegram"
	"go-todobot/internal/domain/events"
	"log"

	"github.com/segmentio/kafka-go"
)

type Fetcher struct {
	reader *kafka.Reader
}

func NewKafkaFetcher(brokers []string, topic, groupID string) *Fetcher {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &Fetcher{
		r,
	}
}

func (f *Fetcher) Fetch(limit int) ([]events.Event, error) {
	log.Println("Attempting to fetch from Kafka...")

	msgs, err := f.reader.FetchMessage(context.Background())
	if err != nil {
		log.Printf("Error fetching messages: %v", err)
		return nil, err
	}

	log.Printf("Fetched message from Kafka: %s", string(msgs.Value))

	var event events.Event
	event, err = convertKafkaMessageToEvent(msgs)
	if err != nil {
		return nil, err
	}

	return []events.Event{event}, nil
}

func convertKafkaMessageToEvent(msg kafka.Message) (events.Event, error) {
	log.Println("converting kafka message...")
	var event events.Event

	event.Text = string(msg.Value)

	event.Meta = telegram.Meta{ChatID: 430746829, Username: "jugglyyfe"}
	event.Type = events.Message

	log.Printf("Converted Kafka message to event: %v", event)

	return event, nil
}

func (f *Fetcher) Close() error {
	return f.reader.Close()
}
