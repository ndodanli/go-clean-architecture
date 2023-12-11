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
)

type ConfirmForgotPasswordCodeQueryHandler struct {
	UOW         uow.IUnitOfWork
	AppServices *services.AppServices
	Logger      logger.ILogger
	TM          *postgresql.TxSessionManager
}

type ConfirmForgotPasswordCodeQuery struct {
	Email              string `body:"email" validate:"required,email"`
	Code               string `body:"code" validate:"required"`
	NewPassword        string `body:"newPassword" validate:"required"`
	NewPasswordConfirm string `body:"newPasswordConfirm" validate:"required"`
}

type ConfirmForgotPasswordCodeQueryResponse struct {
}

func (h *ConfirmForgotPasswordCodeQueryHandler) Handle(echoCtx echo.Context, query *ConfirmForgotPasswordCodeQuery) *baseres.Result[*ConfirmForgotPasswordCodeQueryResponse, error, struct{}] {
	result := baseres.NewResult[*ConfirmForgotPasswordCodeQueryResponse, error, struct{}]()
	ctx := echoCtx.Request().Context()
	appUserRepo := h.UOW.AppUserRepo(ctx, h.TM)

	appUser, err := appUserRepo.FindOneByEmail(query.Email, []string{"id", "password", "email_confirmed", "fp_email_confirmation_details"})
	if err != nil {
		return result.Err(err)
	}

	if appUser == nil {
		return result.Err(httperr.AppUserNotFoundError)
	}

	if !appUser.EmailConfirmed {
		return result.Err(httperr.CannotChangePasswordEmailNotConfirmedError)
	}

	//if time.Now().After(appUser.FpEmailConfirmationDetails.ExpiresAt) {
	//	return result.Err(httperr.ConfirmationCodeExpiredError)
	//}

	if appUser.FpEmailConfirmationDetails.Code != query.Code {
		return result.Err(httperr.InvalidConfirmationCodeError)
	}

	if query.NewPassword != query.NewPasswordConfirm {
		return result.Err(httperr.PasswordsDoNotMatchError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(appUser.Password), []byte(query.NewPassword))
	if err != nil {
		if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return result.Err(err)
		} else {
			return result.Err(httperr.PasswordCannotBeSameAsOldError)
		}
	}

	appUser.FpEmailConfirmationDetails.Code = ""

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(query.NewPassword), bcrypt.DefaultCost)
	_, err = appUserRepo.PatchAppUser(appUser.Id, map[string]interface{}{
		"password":                      newPasswordHash,
		"fp_email_confirmation_details": appUser.FpEmailConfirmationDetails,
	})
	if err != nil {
		return result.Err(err)
	}

	return result.Ok()
}
