package inmem

import (
	"context"

	"github.com/VasySS/avito-winter-2025/internal/entity"
)

func (r *Repository) SetUser(_ context.Context, username string, user entity.User) {
	r.cache.Set(username, user, 1)
}

func (r *Repository) GetUser(_ context.Context, username string) (entity.User, bool) {
	user, ok := r.cache.Get(username)
	if !ok {
		return entity.User{}, false
	}

	if user == nil {
		return entity.User{}, false
	}

	userVal, ok := user.(entity.User)
	if !ok {
		return entity.User{}, false
	}

	return userVal, true
}
