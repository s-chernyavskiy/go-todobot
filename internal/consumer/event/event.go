package event

import (
	"fmt"
	"go-todobot/internal/domain/events"
	"log/slog"
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
	slog.Info("Starting consumer...")
	for {
		event, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			slog.Error("consumer error:", err.Error())

			continue
		}

		slog.Info(fmt.Sprintf("Fetched %d events: ", len(event)))

		if len(event) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(event); err != nil {
			slog.Error("failed to handle events: ", err)

			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	slog.Info("handling events...")

	for _, event := range events {
		slog.Info("got event:", event.Text)

		if err := c.processor.Process(event); err != nil {
			slog.Error("cant handle event:", err.Error())

			continue
		}
	}

	return nil
}
