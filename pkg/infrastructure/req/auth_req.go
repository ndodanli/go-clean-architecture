package req

import "github.com/google/uuid"

type LoginRequest struct {
	UserID   string `json:"-"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type RefreshTokenRequest struct {
	RefreshToken uuid.UUID `param:"refreshToken" validate:"required"`
}
