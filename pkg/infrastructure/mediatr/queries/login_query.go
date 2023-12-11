package queries

import (
	"errors"
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type LoginQueryHandler struct {
	UOW         uow.IUnitOfWork
	AppServices *services.AppServices
	Logger      logger.ILogger
	TM          *postgresql.TxSessionManager
}

type LoginQuery struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type LoginQueryResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *LoginQueryHandler) Handle(echoCtx echo.Context, query *LoginQuery) *baseres.Result[*LoginQueryResponse, error, struct{}] {
	result := baseres.NewResult[*LoginQueryResponse, error, struct{}]()
	ctx := echoCtx.Request().Context()
	authRepo := h.UOW.AuthRepo(ctx, h.TM)
	repoRes, err := authRepo.GetIdAndPasswordWithUsername(query.Username)

	if err != nil {
		return result.Err(err)
	}

	if repoRes == nil {
		return result.Err(httperr.UsernameOrPasswordIncorrectError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(repoRes.Password), []byte(query.Password))
	if err != nil {
		if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return result.Err(err)
		} else {
			return result.Err(httperr.PasswordCannotBeSameAsOldError)
		}
	}

	var accessToken string
	accessToken, err = h.AppServices.JWTService.GenerateAccessToken(strconv.FormatInt(repoRes.ID, 10))

	if err != nil {
		return result.Err(err)
	}

	refreshToken, expiresAt := h.AppServices.JWTService.GenerateRefreshToken()

	_, err = authRepo.UpsertRefreshToken(repoRes.ID, expiresAt, refreshToken)
	if err != nil {
		return result.Err(err)
	}

	result.Data = &LoginQueryResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.String(),
	}

	return result.Ok()
}
