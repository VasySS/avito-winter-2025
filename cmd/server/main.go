package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/VasySS/avito-winter-2025/internal/app"
	"github.com/VasySS/avito-winter-2025/internal/config"
)

func main() {
	cfg := config.MustInit()

	slog.Info("starting app...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, cfg); err != nil {
		slog.Error("error running app", slog.Any("error", err))
	}
}
