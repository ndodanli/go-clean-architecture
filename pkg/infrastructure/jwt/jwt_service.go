package jwtsvc

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/ndodanli/go-clean-architecture/configs"
	"strings"
	"time"
)

type AuthUser struct {
	ID string
}

type JwtServiceInterface interface {
	GenerateAccessToken(id string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type JwtService struct {
	duration time.Duration
	audience []string
	issuer   string
	secret   []byte
}

func NewJwtService(ac configs.Auth) JwtServiceInterface {
	return &JwtService{
		duration: time.Second * time.Duration(ac.JWT_EXPIRATION_IN_SECONDS),
		audience: strings.Split(ac.JWT_AUDIENCES, ","),
		issuer:   ac.JWT_ISSUER,
		secret:   []byte(ac.JWT_SECRET),
	}
}

func (js *JwtService) GenerateAccessToken(id string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    js.issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Audience:  jwt.ClaimStrings(js.audience),
		Subject:   id,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return t.SignedString(js.secret)
}

func (js *JwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return js.secret, nil
	})
}

// TODO: Add refresh token generation and validation
