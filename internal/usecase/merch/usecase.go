package merch

import (
	"context"

	"github.com/VasySS/avito-winter-2025/internal/dto"
	"github.com/VasySS/avito-winter-2025/internal/entity"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=Repository
type Repository interface {
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	SendCoins(ctx context.Context, req entity.UserTransfer) error
	GetMerch(ctx context.Context, name string) (entity.Merch, error)
	BuyMerch(ctx context.Context, req entity.MerchPurchase) error
	Info(ctx context.Context, userID int64) (dto.InfoResponse, error)
}

type Usecase struct {
	repo Repository
}

func New(repo Repository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}
