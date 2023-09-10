package services

import (
	"crypto/rand"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"math/big"
)

type AuthServiceInterface interface {
}

type AuthService struct {
	uow uow.UnitOfWorkInterface
}

func NewAuthService(uow uow.UnitOfWorkInterface) *AuthService {
	return &AuthService{uow: uow}
}

func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result string
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		result += string(charset[randomIndex.Int64()])
	}

	return result, nil
}
