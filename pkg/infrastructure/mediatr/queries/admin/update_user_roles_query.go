package adminqueries

import (
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/pg"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/pg/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"go/types"
)

type UpdateUserRolesQueryHandler struct {
	UOW         uow.IUnitOfWork
	AppServices *services.AppServices
	Logger      logger.ILogger
	TM          *pg.TxSessionManager
}

type UpdateUserRolesQuery struct {
	UserId  int64   `query:"userId" json:"userId" validate:"required"`
	RoleIds []int64 `query:"roleIds" json:"roleIds" validate:"required"`
}

type UpdateUserRolesQueryResponse struct {
}

func (h *UpdateUserRolesQueryHandler) Handle(echoCtx echo.Context, query *UpdateUserRolesQuery) *baseres.Result[*UpdateUserRolesQueryResponse, error, struct{}] {
	result := baseres.NewResult[*UpdateUserRolesQueryResponse, error, struct{}](&UpdateUserRolesQueryResponse{})
	ctx := echoCtx.Request().Context()

	_, err := pg.ExecDefaultTx(ctx, h.TM, func(tx pgx.Tx) (*types.Nil, error) {
		isAppUserExist, err := h.UOW.AppUserRepo(ctx, h.TM).IsUserExist(query.UserId)
		if err != nil {
			return nil, err
		}
		if !isAppUserExist {
			return nil, httperr.AppUserNotFoundError
		}
		qs := pg.NewQueryString(`UPDATE app_user`).
			AddSet("AND", "role_ids", query.RoleIds).
			AddWhere("AND", "id", "=", query.UserId)

		_, err = tx.Exec(ctx, qs.String(), qs.Args()...)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		return result.Err(err)
	}

	return result.Ok()
}
