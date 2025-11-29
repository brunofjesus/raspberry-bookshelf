package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/brunofjesus/raspberry-bookshelf/internal/service"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(log)

	app := service.New()
	if err := app.Run(context.Background()); err != nil {
		panic(err)
	}
}
