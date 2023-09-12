package jwtsvc

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/ndodanli/go-clean-architecture/configs"
	"strings"
	"time"
)

type AuthUser struct {
	ID int64
}

type JWTServiceInterface interface {
	GenerateAccessToken(id string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	GenerateRefreshToken(id string) (string, error)
}

type JWTService struct {
	duration time.Duration
	audience []string
	issuer   string
	secret   []byte
}

func NewJWTService(ac configs.Auth) JWTServiceInterface {
	return &JWTService{
		duration: time.Second * time.Duration(ac.JWT_EXPIRATION_IN_SECONDS),
		audience: strings.Split(ac.JWT_AUDIENCES, ","),
		issuer:   ac.JWT_ISSUER,
		secret:   []byte(ac.JWT_SECRET),
	}
}

func (js *JWTService) GenerateAccessToken(id string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    js.issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(999999999))),
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

func (js *JWTService) GenerateRefreshToken(id string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject: id,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(js.secret)
}

// TODO: Add refresh token generation and validation
