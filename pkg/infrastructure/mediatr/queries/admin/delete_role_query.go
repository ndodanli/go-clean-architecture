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
	"github.com/ndodanli/go-clean-architecture/pkg/utils/pgutils"
)

type DeleteRoleQueryHandler struct {
	UOW         uow.IUnitOfWork
	AppServices *services.AppServices
	Logger      logger.ILogger
	TM          *pg.TxSessionManager
}

type DeleteRoleQuery struct {
	RoleId int64 `param:"roleId" validate:"required"`
}

type DeleteRoleQueryResponse struct {
	RoleId int64 `json:"roleId"`
}

func (h *DeleteRoleQueryHandler) Handle(echoCtx echo.Context, query *DeleteRoleQuery) *baseres.Result[*DeleteRoleQueryResponse, error, struct{}] {
	result := baseres.NewResult[*DeleteRoleQueryResponse, error, struct{}](&DeleteRoleQueryResponse{})
	ctx := echoCtx.Request().Context()

	roleId, err := pg.ExecDefaultTx(ctx, h.TM, func(tx pgx.Tx) (int64, error) {
		var roleIdStruct []struct {
			Id int64 `json:"id"`
		}
		// Soft delete role
		rows, err := tx.Query(ctx, `UPDATE role SET deleted_at = NOW() WHERE id = $1 AND deleted_at = '0001-01-01 00:00:00+00' RETURNING id`, query.RoleId)
		err = pgutils.ScanRowsToStructs(
			rows,
			&roleIdStruct,
		)
		if err != nil {
			return -1, err
		}

		if len(roleIdStruct) == 0 {
			return -1, httperr.NotFoundError("Role")
		}

		return roleIdStruct[0].Id, nil
	})
	if err != nil {
		return result.Err(err)
	}

	result.Data.RoleId = roleId

	return result.Ok()
}
