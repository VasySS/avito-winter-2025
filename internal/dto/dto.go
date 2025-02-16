package dto

import "time"

type AuthUser struct {
	CurTime  time.Time
	Username string `json:"username"`
	Password string `json:"password"`
}
