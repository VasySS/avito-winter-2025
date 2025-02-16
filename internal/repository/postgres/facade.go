package postgres

import (
	"context"

	"github.com/VasySS/avito-winter-2025/internal/dto"
	"github.com/VasySS/avito-winter-2025/internal/entity"
)

type Cache interface {
	GetInfo(ctx context.Context, userID int64) (dto.InfoResponse, bool)
	SetInfo(ctx context.Context, userID int64, info dto.InfoResponse)
	GetMerch(ctx context.Context, name string) (entity.Merch, bool)
	SetMerch(ctx context.Context, name string, merch entity.Merch)
	GetUser(ctx context.Context, username string) (entity.User, bool)
	SetUser(ctx context.Context, username string, user entity.User)
}

type Facade struct {
	*Repository
	cache Cache
}

func NewFacade(repo *Repository, cache Cache) *Facade {
	if cache == nil {
		cache = &NoCache{}
	}

	return &Facade{
		Repository: repo,
		cache:      cache,
	}
}

func (f *Facade) Info(ctx context.Context, userID int64) (dto.InfoResponse, error) {
	info, ok := f.cache.GetInfo(ctx, userID)
	if ok {
		return info, nil
	}

	return f.Repository.Info(ctx, userID)
}

func (f *Facade) GetMerch(ctx context.Context, name string) (entity.Merch, error) {
	merch, ok := f.cache.GetMerch(ctx, name)
	if ok {
		return merch, nil
	}

	return f.Repository.GetMerch(ctx, name)
}

func (f *Facade) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	user, ok := f.cache.GetUser(ctx, username)
	if ok {
		return user, nil
	}

	return f.Repository.GetUserByUsername(ctx, username)
}

func (f *Facade) SendCoins(ctx context.Context, req entity.UserTransfer) error {
	return f.txManager.RunRepeatableRead(ctx, func(ctx context.Context) error { //nolint:wrapcheck
		return f.Repository.SendCoins(ctx, req)
	})
}

func (f *Facade) BuyMerch(ctx context.Context, req entity.MerchPurchase) error {
	return f.txManager.RunRepeatableRead(ctx, func(ctx context.Context) error { //nolint:wrapcheck
		return f.Repository.BuyMerch(ctx, req)
	})
}
