package postgres

import "context"

type TransactionManager interface {
	RunReadCommitted(ctx context.Context, fn func(context.Context) error) error
	RunRepeatableRead(ctx context.Context, fn func(context.Context) error) error
	RunSerializable(ctx context.Context, fn func(context.Context) error) error
	GetQueryEngine(ctx context.Context) QueryEngine
}

type Repository struct {
	txManager TransactionManager
}

func New(txManager TransactionManager) *Repository {
	return &Repository{
		txManager: txManager,
	}
}
