package inmem

import (
	"context"
	"strconv"

	"github.com/VasySS/avito-winter-2025/internal/dto"
	"github.com/VasySS/avito-winter-2025/internal/entity"
)

func (r *Repository) SetMerch(ctx context.Context, merchName string, merch entity.Merch) {
	r.cache.Set(merchName, merch, 1)
}

func (r *Repository) GetMerch(ctx context.Context, name string) (entity.Merch, bool) {
	merch, ok := r.cache.Get(name)
	if !ok {
		return entity.Merch{}, false
	}

	if merch == nil {
		return entity.Merch{}, false
	}

	merchVal, ok := merch.(entity.Merch)
	if !ok {
		return entity.Merch{}, false
	}

	return merchVal, true
}

func (r *Repository) SetInfo(ctx context.Context, userID int64, info dto.InfoResponse) {
	r.cache.Set(strconv.FormatInt(userID, 10), info, 1)
}

func (r *Repository) GetInfo(ctx context.Context, userID int64) (dto.InfoResponse, bool) {
	info, ok := r.cache.Get(strconv.FormatInt(userID, 10))
	if !ok {
		return dto.InfoResponse{}, false
	}

	if info == nil {
		return dto.InfoResponse{}, false
	}

	infoVal, ok := info.(dto.InfoResponse)
	if !ok {
		return dto.InfoResponse{}, false
	}

	return infoVal, true
}
