package app

import (
	"context"
	"log/slog"
)

func Run(ctx context.Context) error {
	closer := NewCloser()

	// code ...

	<-ctx.Done()
	slog.Info("stopping...")

	closeCtx, stop := context.WithCancel(context.Background())
	defer stop()

	if err := closer.Close(closeCtx); err != nil {
		return err
	}

	return nil
}
