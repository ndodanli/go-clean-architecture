package queries

import (
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/backend-api/pkg/core/response"
	"github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg"
	uow "github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg/unit_of_work"
	"github.com/ndodanli/backend-api/pkg/infrastructure/services"
	"github.com/ndodanli/backend-api/pkg/logger"
)

type RegisterQueryHandler struct {
	UOW         uow.IUnitOfWork
	AppServices *services.AppServices
	Logger      logger.ILogger
	TM          *pg.TxSessionManager
}

type RegisterQuery struct {
}

type RegisterQueryResponse struct {
}

func (h *RegisterQueryHandler) Handle(echoCtx echo.Context, query *RegisterQuery) *baseres.Result[*RegisterQueryResponse, error, struct{}] {
	result := baseres.NewResult[*RegisterQueryResponse, error, struct{}](nil)
	//ctx := echoCtx.Request().Context()
	//authRepo := h.UOW.AuthRepo(ctx, h.TM)

	return result.Ok()
}
