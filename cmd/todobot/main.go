package main

import (
	"context"
	tg "go-todobot/internal/client/telegram"
	"go-todobot/internal/config"
	"go-todobot/internal/consumer/event"
	"go-todobot/internal/storage/postgresql"
	"go-todobot/internal/telegram"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	c := config.MustLoad()

	log := initLogger(c.Env)

	tgClient := telegram.New(c.TgBotHost, c.Token)

	conn, err := pgxpool.New(context.Background(), c.DatabaseURL)
	if err != nil {
		log.Error("failed to establish connection to Postgresql:", err)
	}
	db := postgresql.New(conn)

	processor := tg.New(tgClient, db)

	consumer := event.New(processor, processor, c.BatchSize)

	if err := consumer.Start(); err != nil {
		log.Info("consumer stopped: ", err)
	}
}

func initLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
