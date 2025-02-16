package postgres

import (
	"context"
	"fmt"

	"github.com/VasySS/avito-winter-2025/internal/entity"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) SendCoins(ctx context.Context, req entity.UserTransfer) error {
	tx := r.txManager.GetQueryEngine(ctx)

	receiverQuery := `
		UPDATE user_info
		SET balance = balance + @amount
		WHERE id = @id
	`

	_, err := tx.Exec(ctx, receiverQuery, pgx.NamedArgs{
		"id":     req.ReceiverUserID,
		"amount": req.Amount,
	})
	if err != nil {
		return fmt.Errorf("failed to save received coins: %w", err)
	}

	senderQuery := `
		UPDATE user_info
		SET balance = balance - @amount
		WHERE id = @id
	`

	_, err = tx.Exec(ctx, senderQuery, pgx.NamedArgs{
		"id":     req.SenderUserID,
		"amount": req.Amount,
	})
	if err != nil {
		return fmt.Errorf("failed to save sent coins: %w", err)
	}

	transferQuery := `
		INSERT INTO user_transfer (sender_user_id, receiver_user_id, amount, created_at)
		VALUES (@sender_user_id, @receiver_user_id, @amount, @created_at)
	`

	_, err = tx.Exec(ctx, transferQuery, pgx.NamedArgs{
		"sender_user_id":   req.SenderUserID,
		"receiver_user_id": req.ReceiverUserID,
		"amount":           req.Amount,
		"created_at":       req.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to create transfer: %w", err)
	}

	return nil
}

func (r *Repository) GetMerch(ctx context.Context, name string) (entity.Merch, error) {
	tx := r.txManager.GetQueryEngine(ctx)

	query := `
		SELECT *
		FROM merch_item
		WHERE name = @name
	`

	var merch entity.Merch

	err := pgxscan.Get(ctx, tx, &merch, query, pgx.NamedArgs{"name": name})
	if pgxscan.NotFound(err) {
		return entity.Merch{}, entity.ErrMerchItemNotFound
	} else if err != nil {
		return entity.Merch{}, fmt.Errorf("failed to get merch: %w", err)
	}

	return merch, nil
}

func (r *Repository) BuyMerch(ctx context.Context, req entity.MerchPurchase) error {
	tx := r.txManager.GetQueryEngine(ctx)

	userQuery := `
		UPDATE user_info
		SET balance = balance - @price
		WHERE id = @id
	`

	_, err := tx.Exec(ctx, userQuery, pgx.NamedArgs{
		"id":    req.UserID,
		"price": req.Price,
	})
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	merchQuery := `
		INSERT INTO merch_purchase (user_id, merch_item_id, price, created_at)
		VALUES (@user_id, @merch_item_id, @price, @created_at)
	`

	_, err = tx.Exec(ctx, merchQuery, pgx.NamedArgs{
		"user_id":       req.UserID,
		"merch_item_id": req.MerchItemID,
		"price":         req.Price,
		"created_at":    req.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to create merch purchase: %w", err)
	}

	return nil
}
