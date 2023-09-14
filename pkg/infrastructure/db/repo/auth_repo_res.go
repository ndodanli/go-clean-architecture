package repo

import "time"

type GetOnlyIdRepoRes struct {
	ID       int64  `json:"id"`
	Password string `json:"password"`
}

type RefreshTokenRepoRes struct {
	Token     string    `json:"token"`
	AppUserId int64     `json:"app_user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type RefreshTokenWithUUIDRepoRes struct {
	ID        int64     `json:"id"`
	AppUserId int64     `json:"app_user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}
