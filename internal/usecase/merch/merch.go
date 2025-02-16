package merch

import (
	"context"

	"github.com/VasySS/avito-winter-2025/internal/dto"
	"github.com/VasySS/avito-winter-2025/internal/entity"
)

func (u *Usecase) BuyItem(ctx context.Context, item string) error {
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

func (u *Usecase) Info(ctx context.Context) error {
	return nil
}
