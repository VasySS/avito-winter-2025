package entity

import (
	"errors"
	"time"
)

var ErrMerchItemNotFound = errors.New("merch item not found")

type Merch struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

type MerchPurchase struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userID"`
	MerchItemID int64     `json:"merchItemID"`
	Price       int       `json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
}
