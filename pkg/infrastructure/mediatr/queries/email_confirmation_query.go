package queries

import (
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/backend-api/pkg/core/response"
	httperr "github.com/ndodanli/backend-api/pkg/errors"
	"github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg"
	entity "github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg/entity/app_user"
	uow "github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg/unit_of_work"
	"github.com/ndodanli/backend-api/pkg/infrastructure/services"
	"github.com/ndodanli/backend-api/pkg/logger"
	"github.com/ndodanli/backend-api/pkg/utils"
	"time"
)

type EmailConfirmationQueryHandler struct {
	UOW         uow.IUnitOfWork
	AppServices *services.AppServices
	Logger      logger.ILogger
	TM          *pg.TxSessionManager
}

type EmailConfirmationQuery struct {
	Uid  string `query:"uid" validate:"required"`
	Code string `query:"code" validate:"required"`
}

type EmailConfirmationQueryResponse struct {
}

func (h *EmailConfirmationQueryHandler) Handle(echoCtx echo.Context, query *EmailConfirmationQuery) *baseres.Result[*EmailConfirmationQueryResponse, error, struct{}] {
	result := baseres.NewResult[*EmailConfirmationQueryResponse, error, struct{}](nil)
	ctx := echoCtx.Request().Context()
	appUserRepo := h.UOW.AppUserRepo(ctx, h.TM)

	var int64Uid int64
	var err error
	int64Uid, err = utils.ParseInt64(query.Uid)
	if err != nil {
		return result.Err(err)
	}

	var appUser *entity.AppUser
	appUser, err = appUserRepo.FindOneById(int64Uid, []string{"email_confirmed", "email_confirmation"})
	if err != nil {
		return result.Err(httperr.AppUserNotFoundError)
	}

	if appUser.EmailConfirmed {
		return result.Err(httperr.EmailAlreadyConfirmedError)
	}

	if appUser.EmailConfirmation.ExpiresAt.Before(time.Now()) {
		return result.Err(httperr.ConfirmationCodeExpiredError)
	}

	if appUser.EmailConfirmation.Code != query.Code {
		return result.Err(httperr.InvalidConfirmationCodeError)
	}

	appUser.EmailConfirmation.Code = ""

	_, err = appUserRepo.PatchUser(appUser.Id, map[string]interface{}{
		"email_confirmed":    true,
		"email_confirmation": appUser.EmailConfirmation,
	})
	if err != nil {
		return result.Err(err)
	}

	result.SetMessage("Email doğrulama başarılı")

	return result.Ok()
}
