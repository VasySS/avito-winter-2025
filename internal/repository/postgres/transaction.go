package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txManagerKey struct{}

type QueryEngine interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{
		pool: pool,
	}
}

func (tm *TxManager) RunReadCommitted(ctx context.Context, fn func(context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}

	return tm.beginFunc(ctx, opts, fn)
}

func (tm *TxManager) RunSerializable(ctx context.Context, fn func(context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}

	return tm.beginFunc(ctx, opts, fn)
}

func (tm *TxManager) GetQueryEngine(ctx context.Context) QueryEngine {
	tx, ok := ctx.Value(txManagerKey{}).(QueryEngine)
	if ok && tx != nil {
		return tx
	}

	return tm.pool
}

func (tm *TxManager) beginFunc(
	ctx context.Context,
	txOpts pgx.TxOptions,
	fn func(txCtx context.Context) error,
) error {
	tx, err := tm.pool.BeginTx(ctx, txOpts)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	ctx = context.WithValue(ctx, txManagerKey{}, tx)
	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
