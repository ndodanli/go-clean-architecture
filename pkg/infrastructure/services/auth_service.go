package services

import (
	"context"
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/req"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/res"
	"github.com/ndodanli/go-clean-architecture/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type IAuthService interface {
	Login(ctx context.Context, payload req.LoginRequest, ts *postgresql.TxSessionManager) *baseres.Result[res.LoginRes, error, struct{}]
	RefreshToken(ctx context.Context, payload req.RefreshTokenRequest, manager *postgresql.TxSessionManager) *baseres.Result[res.RefreshTokenRes, error, struct{}]
}

type AuthService struct {
	uow        uow.UnitOfWorkInterface
	jwtService IJWTService
}

func NewAuthService(uow uow.UnitOfWorkInterface, jwtService IJWTService) IAuthService {
	return &AuthService{uow: uow, jwtService: jwtService}
}

func (s *AuthService) Login(ctx context.Context, payload req.LoginRequest, ts *postgresql.TxSessionManager) *baseres.Result[res.LoginRes, error, struct{}] {
	result := baseres.NewResult[res.LoginRes, error, struct{}]()

	authRepo := s.uow.AuthRepo(ctx)

	repoRes, err := authRepo.GetIdAndPasswordWithUsername(payload.Username, ts)

	if err != nil {
		return result.Err(err)
	}

	if repoRes == nil {
		return result.Err(httperr.UsernameOrPasswordIncorrectError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(repoRes.Password), []byte(payload.Password))
	if err != nil {
		return result.Err(httperr.UsernameOrPasswordIncorrectError)
	}

	var accessToken string
	accessToken, err = s.jwtService.GenerateAccessToken(strconv.FormatInt(repoRes.ID, 10))

	if err != nil {
		return result.Err(err)
	}

	refreshToken, expiresAt := s.jwtService.GenerateRefreshToken()

	_, err = authRepo.UpsertRefreshToken(repoRes.ID, expiresAt, refreshToken, ts)
	if err != nil {
		return result.Err(err)
	}

	result.Data = res.LoginRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.String(),
	}

	return result.Ok()
}

func (s *AuthService) RefreshToken(ctx context.Context, payload req.RefreshTokenRequest, ts *postgresql.TxSessionManager) *baseres.Result[res.RefreshTokenRes, error, struct{}] {
	result := baseres.NewResult[res.RefreshTokenRes, error, struct{}]()

	authRepo := s.uow.AuthRepo(ctx)

	repoRes, err := authRepo.GetRefreshTokenWithUUID(payload.RefreshToken, ts)
	if err != nil {
		return result.Err(err)
	}

	if repoRes == nil {
		return result.Err(httperr.RefreshTokenNotFoundError)
	}

	if repoRes.ExpiresAt.Before(utils.UTCNow()) {
		return result.Err(httperr.RefreshTokenExpiredError)
	}

	refreshToken, expiresAt := s.jwtService.GenerateRefreshToken()
	var accessToken string
	accessToken, err = s.jwtService.GenerateAccessToken(strconv.FormatInt(repoRes.AppUserId, 10))
	if err != nil {
		return result.Err(err)
	}

	_, err = authRepo.UpdateRefreshToken(repoRes.ID, expiresAt, refreshToken, ts)
	if err != nil {
		return result.Err(err)
	}

	result.Data = res.RefreshTokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.String(),
	}
	return result.Ok()
}
