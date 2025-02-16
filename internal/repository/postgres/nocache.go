package postgres

import (
	"context"

	"github.com/VasySS/avito-winter-2025/internal/dto"
	"github.com/VasySS/avito-winter-2025/internal/entity"
)

type NoCache struct{}

func (r *NoCache) Close() {}

func (r *NoCache) SetUser(_ context.Context, _ string, _ entity.User) {}

func (r *NoCache) GetUser(_ context.Context, _ string) (entity.User, bool) {
	return entity.User{}, false
}

func (r *NoCache) SetMerch(_ context.Context, _ string, _ entity.Merch) {}

func (r *NoCache) GetMerch(_ context.Context, _ string) (entity.Merch, bool) {
	return entity.Merch{}, false
}

func (r *NoCache) Info(_ context.Context, _ int64) (dto.InfoResponse, bool) {
	return dto.InfoResponse{}, false
}

func (r *NoCache) SetInfo(_ context.Context, _ int64, _ dto.InfoResponse) {}

func (r *NoCache) GetInfo(_ context.Context, _ int64) (dto.InfoResponse, bool) {
	return dto.InfoResponse{}, false
}
