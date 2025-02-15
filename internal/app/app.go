package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"avito-task/internal/config"
	httpRouter "avito-task/internal/controller/http"
	"avito-task/internal/repository/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(ctx context.Context, cfg config.Config) error {
	closer := NewCloser()

	pg, err := newPostgres(ctx, closer, cfg.DBConnURL)
	if err != nil {
		return err
	}

	_ = pg

	r := httpRouter.NewRouter()

	go startHTTP(r, closer)

	<-ctx.Done()
	slog.Info("stopping...")

	closeCtx, stop := context.WithCancel(context.Background())
	defer stop()

	if err := closer.Close(closeCtx); err != nil {
		return err
	}

	return nil
}

func startHTTP(r http.Handler, closer *Closer) {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	closer.AddWithCtx(server.Shutdown)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("could not start http server: %s", err.Error())
	}
}

func newPostgres(ctx context.Context, closer *Closer, connURL string) (*postgres.Facade, error) {
	pool, err := pgxpool.New(ctx, connURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	closer.Add(pool.Close)

	txManager := postgres.NewTxManager(pool)
	pgStorage := postgres.New(txManager)
	pgFacade := postgres.NewFacade(pgStorage)

	return pgFacade, nil
}
