package adminqueries

import (
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/backend-api/pkg/core/response"
	httperr "github.com/ndodanli/backend-api/pkg/errors"
	"github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg"
	uow "github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg/unit_of_work"
	"github.com/ndodanli/backend-api/pkg/infrastructure/services"
	"github.com/ndodanli/backend-api/pkg/logger"
	"github.com/ndodanli/backend-api/pkg/utils/pgutils"
)

type AddOrUpdateRoleQueryHandler struct {
	UOW         uow.IUnitOfWork
	AppServices *services.AppServices
	Logger      logger.ILogger
	TM          *pg.TxSessionManager
}

type AddOrUpdateRoleQuery struct {
	RoleName    string  `json:"roleName" validate:"required"`
	Description string  `json:"description"`
	EndpointIds []int64 `json:"endpointIds" validate:"required"`
}

type AddOrUpdateRoleQueryResponse struct {
	RoleId int64 `json:"roleId"`
}

func (h *AddOrUpdateRoleQueryHandler) Handle(echoCtx echo.Context, query *AddOrUpdateRoleQuery) *baseres.Result[*AddOrUpdateRoleQueryResponse, error, struct{}] {
	result := baseres.NewResult[*AddOrUpdateRoleQueryResponse, error, struct{}](&AddOrUpdateRoleQueryResponse{})
	ctx := echoCtx.Request().Context()

	roleId, err := pg.ExecDefaultTx(ctx, h.TM, func(tx pgx.Tx) (int64, error) {
		// check if endpoint ids are valid
		var endpointIdsStruct []pg.IdStruct
		rows, err := tx.Query(ctx, `SELECT id FROM endpoint WHERE id = ANY($1)`, query.EndpointIds)
		err = pgutils.ScanRowsToStructs(
			rows,
			&endpointIdsStruct,
		)
		if err != nil {
			return -1, err
		}

		if endpointIdsStruct == nil || len(endpointIdsStruct) == 0 || len(endpointIdsStruct) != len(query.EndpointIds) {
			return -1, httperr.EndpointIdsAreNotValid
		}

		var roleIdStruct []pg.IdStruct
		rows, err = tx.Query(ctx, `INSERT INTO role (name, description,endpoint_ids) 
									VALUES ($1, $2, $3)
									ON CONFLICT (name) DO UPDATE
     								SET description = $2,
         							endpoint_ids = $3
									RETURNING id
									`, query.RoleName, query.Description, query.EndpointIds)
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
