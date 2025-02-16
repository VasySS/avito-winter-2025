package postgres

import (
	"context"
	"fmt"

	"github.com/VasySS/avito-winter-2025/internal/entity"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateUser(ctx context.Context, user entity.User) error {
	tx := r.txManager.GetQueryEngine(ctx)

	query := `
		INSERT INTO user_info (username, password, created_at)
		VALUES (@username, @password, @created_at)
	`

	_, err := tx.Exec(ctx, query, pgx.NamedArgs{
		"username":   user.Username,
		"password":   user.Password,
		"created_at": user.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	tx := r.txManager.GetQueryEngine(ctx)

	query := `
		SELECT *
		FROM user_info
		WHERE username = @username
	`

	var user entity.User

	err := pgxscan.Get(ctx, tx, &user, query, pgx.NamedArgs{"username": username})
	if pgxscan.NotFound(err) {
		return entity.User{}, entity.ErrUserNotFound
	} else if err != nil {
		return entity.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
