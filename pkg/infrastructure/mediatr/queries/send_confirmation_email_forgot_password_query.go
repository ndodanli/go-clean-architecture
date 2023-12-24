package queries

import (
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/backend-api/pkg/core/response"
	httperr "github.com/ndodanli/backend-api/pkg/errors"
	"github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg"
	uow "github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg/unit_of_work"
	"github.com/ndodanli/backend-api/pkg/infrastructure/services"
	"github.com/ndodanli/backend-api/pkg/logger"
	"github.com/ndodanli/backend-api/pkg/utils"
	"time"
)

type SendConfirmationEmailForgotPasswordQueryHandler struct {
	UOW         uow.IUnitOfWork
	AppServices *services.AppServices
	Logger      logger.ILogger
	TM          *pg.TxSessionManager
}

type SendConfirmationEmailForgotPasswordQuery struct {
	Email string `param:"email" validate:"required,email"`
}

type SendConfirmationEmailForgotPasswordQueryResponse struct {
	Email string `json:"email"`
}

func (h *SendConfirmationEmailForgotPasswordQueryHandler) Handle(echoCtx echo.Context, query *SendConfirmationEmailForgotPasswordQuery) *baseres.Result[*SendConfirmationEmailForgotPasswordQueryResponse, error, struct{}] {
	result := baseres.NewResult[*SendConfirmationEmailForgotPasswordQueryResponse, error, struct{}](nil)
	ctx := echoCtx.Request().Context()
	appUserRepo := h.UOW.AppUserRepo(ctx, h.TM)

	appUser, err := appUserRepo.FindOneByEmail(query.Email, []string{"id", "email_confirmed", "fp_email_confirmation"})
	if err != nil {
		return result.Err(err)
	}

	if appUser == nil {
		return result.Err(httperr.AppUserNotFoundError)
	}

	if !appUser.EmailConfirmed {
		return result.Err(httperr.CannotChangePasswordEmailNotConfirmedError)
	}

	validDate := appUser.FpEmailConfirmation.ExpiresAt.Add(-9 * time.Minute)
	if time.Now().Before(validDate) {
		return result.Err(httperr.CodeRecentlySentError)
	}

	appUser.FpEmailConfirmation.Code = utils.GenerateCodeOnlyNumbers(6)
	appUser.FpEmailConfirmation.ExpiresAt = time.Now().Add(10 * time.Minute)

	err = h.AppServices.SendgridService.SendEmail(query.Email, "Forgot Password", "Your code is: "+appUser.FpEmailConfirmation.Code)
	if err != nil {
		return result.Err(err)
	}

	_, err = appUserRepo.PatchUser(appUser.Id, map[string]interface{}{
		"fp_email_confirmation": appUser.FpEmailConfirmation,
	})
	if err != nil {
		return result.Err(err)
	}

	result.SetMessage("Confirmation email sent")

	return result.Ok()
}
