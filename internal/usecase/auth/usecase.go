package auth

import (
	"context"
	"time"

	"github.com/VasySS/avito-winter-2025/internal/entity"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=Repository
type Repository interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=PasswordHasher
type PasswordHasher interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=TokenGenerator
type TokenGenerator interface {
	NewAccessToken(user entity.User, currentTime time.Time) (string, error)
}

type Usecase struct {
	repo      Repository
	hasher    PasswordHasher
	generator TokenGenerator
}

func New(
	repo Repository,
	hasher PasswordHasher,
	generator TokenGenerator,
) *Usecase {
	return &Usecase{
		repo:      repo,
		hasher:    hasher,
		generator: generator,
	}
}
