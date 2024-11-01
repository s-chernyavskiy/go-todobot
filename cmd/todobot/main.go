package main

import (
	"context"
	telegram2 "go-todobot/internal/client/telegram"
	"go-todobot/internal/config"
	"go-todobot/internal/consumer/event"
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
	db := postgresql.New(conn)

	processor := telegram2.New(tgClient, db)

	consumer := event.New(processor, processor, c.BatchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("consumer stopped: ", err)
	}
}
