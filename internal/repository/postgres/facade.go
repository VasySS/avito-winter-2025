package postgres

import (
	"context"

	"github.com/VasySS/avito-winter-2025/internal/entity"
)

type Facade struct {
	*Repository
}

func NewFacade(repo *Repository) *Facade {
	return &Facade{
		Repository: repo,
	}
}

func (f *Facade) SendCoins(ctx context.Context, req entity.UserTransfer) error {
	return f.txManager.RunRepeatableRead(ctx, func(ctx context.Context) error {
		return f.Repository.SendCoins(ctx, req)
	})
}

func (f *Facade) BuyMerch(ctx context.Context, req entity.MerchPurchase) error {
	return f.txManager.RunRepeatableRead(ctx, func(ctx context.Context) error {
		return f.Repository.BuyMerch(ctx, req)
	})
}
