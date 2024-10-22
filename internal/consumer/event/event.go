package event

import (
	"go-todobot/internal/domain/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	log.Println("Starting consumer...")
	for {
		event, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("consumer error: %v", err.Error())

			continue
		}

		log.Printf("Fetched %d events", len(event))

		if len(event) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(event); err != nil {
			log.Print(err)

			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	log.Println("handling events...")

	for _, event := range events {
		log.Printf("got event: %v", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("cant handle event: %v", err.Error())

			continue
		}
	}

	return nil
}
