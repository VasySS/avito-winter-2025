package postgres

import (
	"context"

	"github.com/VasySS/avito-winter-2025/internal/dto"
	"github.com/VasySS/avito-winter-2025/internal/entity"
)

type NoCache struct{}

func (r *NoCache) Close() {}

func (r *NoCache) SetUser(ctx context.Context, username string, user entity.User) {}

func (r *NoCache) GetUser(ctx context.Context, username string) (entity.User, bool) {
	return entity.User{}, false
}

func (r *NoCache) SetMerch(ctx context.Context, name string, merch entity.Merch) {}

func (r *NoCache) GetMerch(ctx context.Context, name string) (entity.Merch, bool) {
	return entity.Merch{}, false
}

func (r *NoCache) Info(ctx context.Context, userID int64) (dto.InfoResponse, bool) {
	return dto.InfoResponse{}, false
}

func (r *NoCache) SetInfo(ctx context.Context, userID int64, info dto.InfoResponse) {}

func (r *NoCache) GetInfo(ctx context.Context, userID int64) (dto.InfoResponse, bool) {
	return dto.InfoResponse{}, false
}
