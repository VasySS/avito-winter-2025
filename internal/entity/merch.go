package entity

import "time"

type Merch struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}
