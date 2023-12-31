package queries

import (
	"errors"
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/backend-api/pkg/core/response"
	httperr "github.com/ndodanli/backend-api/pkg/errors"
	"github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg"
	uow "github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg/unit_of_work"
	"github.com/ndodanli/backend-api/pkg/infrastructure/services"
	"github.com/ndodanli/backend-api/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type LoginQueryHandler struct {
	UOW         uow.IUnitOfWork
	AppServices *services.AppServices
	Logger      logger.ILogger
	TM          *pg.TxSessionManager
}

type LoginQuery struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginQueryResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *LoginQueryHandler) Handle(echoCtx echo.Context, query *LoginQuery) *baseres.Result[*LoginQueryResponse, error, struct{}] {
	result := baseres.NewResult[*LoginQueryResponse, error, struct{}](nil)
	ctx := echoCtx.Request().Context()
	authRepo := h.UOW.AuthRepo(ctx, h.TM)
	repoRes, err := authRepo.GetIdAndPasswordWithUsername(query.Username)
	if err != nil {
		return result.Err(err)
	}

	if repoRes == nil {
		return result.Err(httperr.UsernameOrPasswordIncorrectError)
	}

	authorizeResponse, err := h.AppServices.JWTService.Authorize(ctx, h.UOW.DB(), repoRes.ID, echoCtx.Path(), echoCtx.Request().Method)
	if err != nil {
		return result.Err(err)
	}

	if !authorizeResponse.IsAuthorized {
		return result.Err(httperr.UnauthorizedError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(repoRes.Password), []byte(query.Password))
	if err != nil {
		if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return result.Err(err)
		} else {
			return result.Err(httperr.UsernameOrPasswordIncorrectError)
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
