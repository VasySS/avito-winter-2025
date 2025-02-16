package merch

import (
	"context"
	"fmt"

	"github.com/VasySS/avito-winter-2025/internal/dto"
	"github.com/VasySS/avito-winter-2025/internal/entity"
)

func (u *Usecase) BuyItem(ctx context.Context, req dto.MerchPurchase) error {
	merch, err := u.repo.GetMerch(ctx, req.MerchName)
	if err != nil {
		return fmt.Errorf("failed to get merch: %w", err)
	}

	user, err := u.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.Balance < merch.Price {
		return entity.ErrLowBalance
	}

	repoReq := entity.MerchPurchase{
		UserID:      user.ID,
		MerchItemID: merch.ID,
		Price:       merch.Price,
		CreatedAt:   req.CurTime,
	}

	if err := u.repo.BuyMerch(ctx, repoReq); err != nil {
		return fmt.Errorf("failed to buy merch: %w", err)
	}

	return nil
}

func (u *Usecase) SendCoin(ctx context.Context, req dto.CoinSend) error {
	if req.FromUser == req.ToUser {
		return entity.ErrSameTransferUser
	}

	senderUser, err := u.repo.GetUserByUsername(ctx, req.FromUser)
	if err != nil {
		return err
	}

	if senderUser.Balance < req.Amount {
		return entity.ErrLowBalance
	}

	receiverUser, err := u.repo.GetUserByUsername(ctx, req.ToUser)
	if err != nil {
		return err
	}

	sendReq := entity.UserTransfer{
		SenderUserID:   senderUser.ID,
		ReceiverUserID: receiverUser.ID,
		Amount:         req.Amount,
		CreatedAt:      req.CurTime,
	}

	if err := u.repo.SendCoins(ctx, sendReq); err != nil {
		return err
	}

	return nil
}

func (u *Usecase) Info(ctx context.Context, username string) (dto.InfoResponse, error) {
	user, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return dto.InfoResponse{}, err
	}

	resp, err := u.repo.Info(ctx, user.ID)
	if err != nil {
		return dto.InfoResponse{}, err
	}

	return resp, nil
}
