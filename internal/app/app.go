package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/VasySS/avito-winter-2025/internal/config"
	httpRouter "github.com/VasySS/avito-winter-2025/internal/controller/http"
	"github.com/VasySS/avito-winter-2025/internal/repository/inmem"
	"github.com/VasySS/avito-winter-2025/internal/repository/postgres"
	"github.com/VasySS/avito-winter-2025/internal/usecase/auth"
	"github.com/VasySS/avito-winter-2025/internal/usecase/merch"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(ctx context.Context, cfg config.Config) error {
	closer := NewCloser()

	pg, err := newPostgres(ctx, closer, cfg)
	if err != nil {
		return err
	}

	authUsecase := auth.New(
		pg,
		auth.NewBcryptPasswordHasher(),
		auth.NewJWTGenerator(cfg.JWTSecret, cfg.AccessTokenTTL),
	)
	merchUsecase := merch.New(pg)

	r := httpRouter.NewRouter(cfg, merchUsecase, authUsecase)

	go startHTTP(r, closer, cfg.ServerPort)

	<-ctx.Done()
	slog.Info("stopping...")

	closeCtx, stop := context.WithCancel(context.Background())
	defer stop()

	if err := closer.Close(closeCtx); err != nil { //nolint:contextcheck
		return err
	}

	return nil
}

func startHTTP(r http.Handler, closer *Closer, port string) {
	server := &http.Server{
		Addr:         ":" + port,
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

func newPostgres(ctx context.Context, closer *Closer, cfg config.Config) (*postgres.Facade, error) {
	connURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
	)

	pool, err := pgxpool.New(ctx, connURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	closer.Add(pool.Close)

	cache, err := inmem.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

	closer.Add(cache.Close)

	txManager := postgres.NewTxManager(pool)
	pgStorage := postgres.New(txManager)
	pgFacade := postgres.NewFacade(pgStorage, cache)

	return pgFacade, nil
}
