package queries

import (
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
)

type RegisterQueryHandler struct {
	UOW         uow.IUnitOfWork
	AppServices *services.AppServices
	Logger      logger.ILogger
	TM          *postgresql.TxSessionManager
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
