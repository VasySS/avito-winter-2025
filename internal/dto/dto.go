package dto

import "time"

type AuthUser struct {
	CurTime  time.Time
	Username string `json:"username"`
	Password string `json:"password"`
}

type CoinSend struct {
	CurTime  time.Time
	FromUser string
	ToUser   string `json:"toUser"`
	Amount   int    `json:"amount"`
}

type MerchPurchase struct {
	CurTime   time.Time
	Username  string
	MerchName string
}
