package models

import "time"

type PasswordResetToken struct {
	ID 	int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Token string `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
	Used bool `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}