package services

import (
	"context"
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	jwtsvc "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/jwt"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/req"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/res"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type AuthServiceInterface interface {
	Login(ctx context.Context, payload req.LoginRequest, ts *postgresql.TxSessionManager) *baseres.Result[res.LoginRes, error, struct{}]
}

type AuthService struct {
	uow        uow.UnitOfWorkInterface
	jwtService jwtsvc.JWTServiceInterface
}

func NewAuthService(uow uow.UnitOfWorkInterface, jwtService jwtsvc.JWTServiceInterface) AuthServiceInterface {
	return &AuthService{uow: uow, jwtService: jwtService}
}

func (s *AuthService) Login(ctx context.Context, payload req.LoginRequest, ts *postgresql.TxSessionManager) *baseres.Result[res.LoginRes, error, struct{}] {
	result := baseres.NewResult[res.LoginRes, error, struct{}]()

	userRepo := s.uow.AppUserRepo(ctx)

	repoRes, err := userRepo.GetIdAndPasswordWithUsername(payload.Username, ts)

	if err != nil {
		return result.Err(err)
	}

	if repoRes == nil {
		return result.Err(httperr.UsernameOrPasswordIncorrect)
	}

	err = bcrypt.CompareHashAndPassword([]byte(repoRes.Password), []byte(payload.Password))
	if err != nil {
		return result.Err(httperr.UsernameOrPasswordIncorrect)
	}

	var accessToken string
	accessToken, err = s.jwtService.GenerateAccessToken(strconv.FormatInt(repoRes.ID, 10))

	if err != nil {
		return result.Err(err)
	}

	var refreshToken string
	refreshToken, err = s.jwtService.GenerateRefreshToken(strconv.FormatInt(repoRes.ID, 10))
	if err != nil {
		return result.Err(err)
	}

	result.Data = res.LoginRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return result.Ok()
}
