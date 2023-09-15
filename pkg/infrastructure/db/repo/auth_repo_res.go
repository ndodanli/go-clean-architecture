package repo

import "time"

type GetIdAndPasswordRepoRes struct {
	ID       int64  `json:"id"`
	Password string `json:"password"`
}

type RefreshTokenWithUUIDRepoRes struct {
	ID        int64     `json:"id"`
	AppUserId int64     `json:"appUserId"`
	ExpiresAt time.Time `json:"expiresAt"`
}
