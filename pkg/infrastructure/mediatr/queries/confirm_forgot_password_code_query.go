package queries

import (
	"errors"
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/pg"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/pg/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type ConfirmForgotPasswordCodeQueryHandler struct {
	UOW         uow.IUnitOfWork
	AppServices *services.AppServices
	Logger      logger.ILogger
	TM          *pg.TxSessionManager
}

type ConfirmForgotPasswordCodeQuery struct {
	Email              string `json:"email" validate:"required,email"`
	Code               string `json:"code" validate:"required"`
	NewPassword        string `json:"newPassword" validate:"required"`
	NewPasswordConfirm string `json:"newPasswordConfirm" validate:"required"`
}

type ConfirmForgotPasswordCodeQueryResponse struct {
}

func (h *ConfirmForgotPasswordCodeQueryHandler) Handle(echoCtx echo.Context, query *ConfirmForgotPasswordCodeQuery) *baseres.Result[*ConfirmForgotPasswordCodeQueryResponse, error, struct{}] {
	result := baseres.NewResult[*ConfirmForgotPasswordCodeQueryResponse, error, struct{}](nil)
	ctx := echoCtx.Request().Context()
	appUserRepo := h.UOW.AppUserRepo(ctx, h.TM)

	appUser, err := appUserRepo.FindOneByEmail(query.Email, []string{"id", "password", "email_confirmed", "fp_email_confirmation"})
	if err != nil {
		return result.Err(err)
	}

	if appUser == nil {
		return result.Err(httperr.AppUserNotFoundError)
	}

	if !appUser.EmailConfirmed {
		return result.Err(httperr.CannotChangePasswordEmailNotConfirmedError)
	}

	//if time.Now().After(appUser.FpEmailConfirmation.ExpiresAt) {
	//	return result.Err(httperr.ConfirmationCodeExpiredError)
	//}

	if appUser.FpEmailConfirmation.Code != query.Code {
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

	appUser.FpEmailConfirmation.Code = ""

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(query.NewPassword), bcrypt.DefaultCost)
	_, err = appUserRepo.PatchAppUser(appUser.Id, map[string]interface{}{
		"password":              newPasswordHash,
		"fp_email_confirmation": appUser.FpEmailConfirmation,
	})
	if err != nil {
		return result.Err(err)
	}

	return result.Ok()
}
