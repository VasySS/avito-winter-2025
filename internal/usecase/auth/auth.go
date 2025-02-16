package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/VasySS/avito-winter-2025/internal/dto"
	"github.com/VasySS/avito-winter-2025/internal/entity"
)

const hashCost = 10

var ErrInvalidPassword = errors.New("invalid password")

func (u *Usecase) AuthUser(ctx context.Context, req dto.AuthUser) (string, error) {
	user, err := u.repo.GetUserByUsername(ctx, req.Username)
	if errors.Is(err, entity.ErrUserNotFound) {
		user, err := u.newUser(ctx, req)
		if err != nil {
			return "", err
		}

		return u.generator.NewAccessToken(user, req.CurTime)
	} else if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if err := u.hasher.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", ErrInvalidPassword
	}

	return u.generator.NewAccessToken(user, req.CurTime)
}

func (u *Usecase) newUser(ctx context.Context, req dto.AuthUser) (entity.User, error) {
	hashedPassword, err := u.hasher.GenerateFromPassword([]byte(req.Password), hashCost)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	userInfo := entity.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		CreatedAt: req.CurTime,
	}

	if err := u.repo.CreateUser(ctx, userInfo); err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	userRepo, err := u.repo.GetUserByUsername(ctx, userInfo.Username)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return userRepo, nil
}
