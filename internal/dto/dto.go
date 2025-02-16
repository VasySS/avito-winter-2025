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

type InventoryItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type CoinTransferReceived struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type CoinTransferSent struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type CoinHistory struct {
	Received []CoinTransferReceived `json:"received"`
	Sent     []CoinTransferSent     `json:"sent"`
}

type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}
