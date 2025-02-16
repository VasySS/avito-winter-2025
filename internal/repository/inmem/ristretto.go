package inmem

import (
	"fmt"

	"github.com/dgraph-io/ristretto/v2"
)

type Repository struct {
	cache *ristretto.Cache[string, any]
}

func New() (*Repository, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, any]{
		NumCounters: 1e5,
		MaxCost:     1 << 32,
		BufferItems: 64,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

	return &Repository{
		cache: cache,
	}, nil
}

func (r *Repository) Close() {
	r.cache.Close()
}
