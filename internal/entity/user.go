package entity

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserTransferNotFound = errors.New("user transfer not found")
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
