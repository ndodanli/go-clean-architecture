package queries

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/backend-api/pkg/core/response"
	httperr "github.com/ndodanli/backend-api/pkg/errors"
	"github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg"
	uow "github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg/unit_of_work"
	"github.com/ndodanli/backend-api/pkg/infrastructure/services"
	"github.com/ndodanli/backend-api/pkg/logger"
	"github.com/ndodanli/backend-api/pkg/utils"
	"strconv"
)

type RefreshTokenQueryHandler struct {
	UOW         uow.IUnitOfWork
	AppServices *services.AppServices
	Logger      logger.ILogger
	TM          *pg.TxSessionManager
}

type RefreshTokenQuery struct {
	RefreshToken uuid.UUID `param:"refreshToken" validate:"required"`
}

type RefreshTokenQueryResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *RefreshTokenQueryHandler) Handle(echoCtx echo.Context, query *RefreshTokenQuery) *baseres.Result[*RefreshTokenQueryResponse, error, struct{}] {
	result := baseres.NewResult[*RefreshTokenQueryResponse, error, struct{}](nil)
	ctx := echoCtx.Request().Context()
	authRepo := h.UOW.AuthRepo(ctx, h.TM)

	repoRes, err := authRepo.GetRefreshTokenWithUUID(query.RefreshToken)
	if err != nil {
		return result.Err(err)
	}

	if repoRes == nil {
		return result.Err(httperr.RefreshTokenNotFoundError)
	}

	if repoRes.ExpiresAt.Before(utils.UTCNow()) {
		return result.Err(httperr.RefreshTokenExpiredError)
	}

	refreshToken, expiresAt := h.AppServices.JWTService.GenerateRefreshToken()
	var accessToken string
	accessToken, err = h.AppServices.JWTService.GenerateAccessToken(strconv.FormatInt(repoRes.AppUserId, 10))
	if err != nil {
		return result.Err(err)
	}

	_, err = authRepo.UpdateRefreshToken(repoRes.ID, expiresAt, refreshToken)
	if err != nil {
		return result.Err(err)
	}

	result.Data = &RefreshTokenQueryResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.String(),
	}

	return result.Ok()
}
