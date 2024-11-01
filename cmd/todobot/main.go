package main

import (
	"context"
	"go-todobot/internal/config"
	"go-todobot/internal/consumer/event"
	"go-todobot/internal/kafka"
	"go-todobot/internal/storage/postgresql"
	"go-todobot/internal/telegram"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	c := config.MustLoad()

	tgClient := telegram.New(c.TgBotHost, c.Token)

	conn, err := pgxpool.New(context.Background(), c.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	kafkaProducer := kafka.NewProducer(tgClient, postgresql.New(conn), []string{c.Kafka.Brokers}, c.Kafka.Topic)
	kafkaConsumer := kafka.NewConsumer([]string{c.Kafka.Brokers}, c.Kafka.Topic, c.Kafka.GroupID, tgClient)

	consumer := event.New(kafkaProducer, kafkaConsumer, c.BatchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("consumer stopped: ", err)
	}

	if err := kafkaConsumer.Close(); err != nil {
		log.Fatal("Error closing Kafka reader: ", err)
	}
}
