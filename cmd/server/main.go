package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"avito-task/internal/app"
	"avito-task/internal/config"
)

func main() {
	cfg := config.MustInit()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, cfg); err != nil {
		slog.Error("error running app", slog.Any("error", err))
	}
}
