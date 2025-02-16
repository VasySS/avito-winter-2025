package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/VasySS/avito-winter-2025/internal/app"
	"github.com/VasySS/avito-winter-2025/internal/config"
)

func main() {
	setupLogger()

	cfg := config.MustInit(".env")

	slog.Info("starting app...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, cfg); err != nil {
		slog.Error("error running app", slog.Any("error", err))
	}
}

func setupLogger() {
	slogLogger := slog.New(
		slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	slog.SetDefault(slogLogger)
}
