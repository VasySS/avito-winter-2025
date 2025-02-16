package entity

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrSameTransferUser = errors.New("can't send coin to the same user")
	ErrLowBalance       = errors.New("not enough coins")
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserTransfer struct {
	ID             int64     `json:"id"`
	SenderUserID   int64     `json:"senderUserID"`
	ReceiverUserID int64     `json:"receiverUserID"`
	Amount         int       `json:"amount"`
	CreatedAt      time.Time `json:"createdAt"`
}
