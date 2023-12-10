package queries

import (
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
)

type SendConfirmationEmailForgotPasswordQueryHandler struct {
	UOW         uow.IUnitOfWork
	AppServices *services.AppServices
	Logger      logger.ILogger
	TM          *postgresql.TxSessionManager
}

type SendConfirmationEmailForgotPasswordQuery struct {
	Email string `param:"email" validate:"required,email"`
}

type SendConfirmationEmailForgotPasswordQueryResponse struct {
	Email string `json:"email"`
}

func (h *SendConfirmationEmailForgotPasswordQueryHandler) Handle(echoCtx echo.Context, query *SendConfirmationEmailForgotPasswordQuery) *baseres.Result[*SendConfirmationEmailForgotPasswordQueryResponse, error, struct{}] {
	result := baseres.NewResult[*SendConfirmationEmailForgotPasswordQueryResponse, error, struct{}]()
	ctx := echoCtx.Request().Context()
	appUserRepo := h.UOW.AppUserRepo(ctx, h.TM)

	appUser, err := appUserRepo.FindOneByEmail(query.Email, []string{})
	if err != nil {
		return result.Err(err)
	}

	if appUser == nil {
		return result.Err(httperr.AppUserNotFoundError)
	}

	return result.Ok()
}
