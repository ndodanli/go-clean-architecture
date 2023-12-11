package queries

import (
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"github.com/ndodanli/go-clean-architecture/pkg/utils"
	"time"
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

	appUser, err := appUserRepo.FindOneByEmail(query.Email, []string{"id", "email_confirmed", "fp_email_confirmation_details"})
	if err != nil {
		return result.Err(err)
	}

	if appUser == nil {
		return result.Err(httperr.AppUserNotFoundError)
	}

	if !appUser.EmailConfirmed {
		return result.Err(httperr.CannotChangePasswordEmailNotConfirmedError)
	}

	validDate := appUser.FpEmailConfirmationDetails.ExpiresAt.Add(-9 * time.Minute)
	if time.Now().Before(validDate) {
		return result.Err(httperr.CodeRecentlySentError)
	}

	appUser.FpEmailConfirmationDetails.Code = utils.GenerateCodeOnlyNumbers(6)
	appUser.FpEmailConfirmationDetails.ExpiresAt = time.Now().Add(10 * time.Minute)

	err = h.AppServices.SendgridService.SendEmail(query.Email, "Forgot Password", "Your code is: "+appUser.FpEmailConfirmationDetails.Code)
	if err != nil {
		return result.Err(err)
	}

	_, err = appUserRepo.PatchAppUser(appUser.Id, map[string]interface{}{
		"fp_email_confirmation_details": appUser.FpEmailConfirmationDetails,
	})
	if err != nil {
		return result.Err(err)
	}

	result.SetMessage("Confirmation email sent")

	return result.Ok()
}
