package postgres

import (
	"context"
	"fmt"

	"github.com/VasySS/avito-winter-2025/internal/entity"
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
