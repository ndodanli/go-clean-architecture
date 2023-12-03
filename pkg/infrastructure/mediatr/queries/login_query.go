package queries

import (
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/constant"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type LoginQueryHandler struct {
	uow        uow.IUnitOfWork
	jwtService services.IJWTService
	logger     logger.ILogger
}

func NewLoginQueryHandler(appServices *services.AppServices, uow uow.IUnitOfWork, logger logger.ILogger) *LoginQueryHandler {
	return &LoginQueryHandler{
		uow:        uow,
		jwtService: appServices.JWTService,
		logger:     logger,
	}
}

type LoginQuery struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type LoginQueryResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *LoginQueryHandler) Handle(echoCtx echo.Context, query *LoginQuery) *baseres.Result[LoginQueryResponse, error, struct{}] {
	result := baseres.NewResult[LoginQueryResponse, error, struct{}]()
	ctx := echoCtx.Request().Context()
	tm := echoCtx.Get(constant.General.TxSessionManagerKey).(*postgresql.TxSessionManager)
	authRepo := h.uow.AuthRepo(ctx)
	repoRes, err := authRepo.GetIdAndPasswordWithUsername(query.Username, tm)

	if err != nil {
		return result.Err(err)
	}

	if repoRes == nil {
		return result.Err(httperr.UsernameOrPasswordIncorrectError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(repoRes.Password), []byte(query.Password))
	if err != nil {
		return result.Err(httperr.UsernameOrPasswordIncorrectError)
	}

	var accessToken string
	accessToken, err = h.jwtService.GenerateAccessToken(strconv.FormatInt(repoRes.ID, 10))

	if err != nil {
		return result.Err(err)
	}

	refreshToken, expiresAt := h.jwtService.GenerateRefreshToken()

	_, err = authRepo.UpsertRefreshToken(repoRes.ID, expiresAt, refreshToken, tm)
	if err != nil {
		return result.Err(err)
	}

	result.Data = LoginQueryResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.String(),
	}

	return result.Ok()
}
