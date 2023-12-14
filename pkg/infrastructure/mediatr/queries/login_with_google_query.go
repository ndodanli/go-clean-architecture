package queries

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/pg"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/pg/unit_of_work"
	oauthcfg "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/oauth_cfg"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"io"
	"net/http"
)

type LoginWithGoogleQueryHandler struct {
	UOW         uow.IUnitOfWork
	AppServices *services.AppServices
	Logger      logger.ILogger
	TM          *pg.TxSessionManager
}

type LoginWithGoogleQuery struct {
	Code string `query:"code" validate:"required"`
}

type LoginWithGoogleQueryResponse struct {
}

func (h *LoginWithGoogleQueryHandler) Handle(echoCtx echo.Context, query *LoginWithGoogleQuery) *baseres.Result[*LoginWithGoogleQueryResponse, error, struct{}] {
	result := baseres.NewResult[*LoginWithGoogleQueryResponse, error, struct{}](nil)
	ctx := echoCtx.Request().Context()
	//authRepo := h.UOW.AuthRepo(ctx, h.TM)
	token, err := oauthcfg.GoogleOauth2Config.Exchange(ctx, query.Code)
	if err != nil {
		return result.Err(err)
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return result.Err(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result.Err(err)
	}

	var googleUserInfo oauthcfg.GoogleUserInfo
	err = json.Unmarshal(body, &googleUserInfo)
	if err != nil {
		return result.Err(err)
	}

	return result.Ok()
}
