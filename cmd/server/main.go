package main

import (
	"context"
	"go_template/internal/app"
	"log/slog"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil {
		slog.Error("error running app", slog.Any("error", err))
	}
}
