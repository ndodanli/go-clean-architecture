package services

import (
	"context"
	"crypto/rand"
	"fmt"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"math/big"
)

type AuthServiceInterface interface {
	GetUser(ctx context.Context) string
	UpdateTestString() string
}

type AuthService struct {
	uow        uow.UnitOfWorkInterface
	testString string
}

func NewAuthService(uow uow.UnitOfWorkInterface) *AuthService {
	str, _ := GenerateRandomString(10)
	return &AuthService{uow: uow, testString: str}
}

func (as *AuthService) GetUser(ctx context.Context) string {
	r := as.uow.UserRepo().TestTx(ctx)
	fmt.Println(r)
	return as.testString
}

func (as *AuthService) UpdateTestString() string {
	str, _ := GenerateRandomString(10)
	as.testString = str

	return as.testString
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
