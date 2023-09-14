package services

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ndodanli/go-clean-architecture/configs"
	"github.com/ndodanli/go-clean-architecture/pkg/utils"
	"strings"
	"time"
)

type AuthUser struct {
	ID int64
}

type IJWTService interface {
	GenerateAccessToken(id string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	GenerateRefreshToken() (u uuid.UUID, expiresAt time.Time)
}

type JWTService struct {
	duration               time.Duration
	refreshTokenExpiration time.Duration
	audience               []string
	issuer                 string
	secret                 []byte
}

func NewJWTService(ac configs.Auth) IJWTService {
	return &JWTService{
		duration:               time.Second * time.Duration(ac.JWT_EXPIRATION_IN_SECONDS),
		refreshTokenExpiration: time.Second * time.Duration(ac.JWT_REFRESH_EXPIRATION_IN_SECONDS),
		audience:               strings.Split(ac.JWT_AUDIENCES, ","),
		issuer:                 ac.JWT_ISSUER,
		secret:                 []byte(ac.JWT_SECRET),
	}
}

func (js *JWTService) GenerateAccessToken(id string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    js.issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Audience:  jwt.ClaimStrings(js.audience),
		Subject:   id,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(js.secret)
}

func (js *JWTService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return js.secret, nil
	})

}

func (js *JWTService) GenerateRefreshToken() (u uuid.UUID, expiresAt time.Time) {
	return uuid.New(), utils.UTCNowAddDuration(js.refreshTokenExpiration)
}
