package adminqueries

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/backend-api/pkg/core/response"
	"github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg"
	entity "github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg/entity/app_user"
	uow "github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg/unit_of_work"
	"github.com/ndodanli/backend-api/pkg/infrastructure/services"
	"github.com/ndodanli/backend-api/pkg/logger"
	"github.com/ndodanli/backend-api/pkg/utils/pgutils"
	"time"
)

type GetUsersQueryHandler struct {
	UOW         uow.IUnitOfWork
	AppServices *services.AppServices
	Logger      logger.ILogger
	TM          *pg.TxSessionManager
}

type GetUsersQuery struct {
	pg.PaginationQuery
	RoleIds     []int64 `query:"roleIds" json:"roleIds"`
	EndpointIds []int64 `query:"endpointIds" json:"endpointIds"`
}

type GetUsersQueryResponse struct {
	AppUsers   []entity.AppUser `json:"appUsers"`
	TotalCount int64            `json:"totalCount"`
}

func (h *GetUsersQueryHandler) Handle(echoCtx echo.Context, query *GetUsersQuery) *baseres.Result[*GetUsersQueryResponse, error, struct{}] {
	result := baseres.NewResult[*GetUsersQueryResponse, error, struct{}](&GetUsersQueryResponse{})
	ctx := echoCtx.Request().Context()
	res, err := pg.ExecDefaultTx(ctx, h.TM, func(tx pgx.Tx) (*GetUsersQueryResponse, error) {
		// assign role ids empty array if not provided
		txRes := GetUsersQueryResponse{}
		timeSt := time.Now()
		qs := pg.NewQueryString(`SELECT * FROM app_user`)

		if len(query.EndpointIds) > 0 {
			// get role ids that have the given endpoint ids
			roleIdsQs := pg.NewQueryString(`SELECT DISTINCT id FROM role`).
				AddWhere("AND", "endpoint_ids", "&&", query.EndpointIds)
			var roleIds []pg.IdStruct
			roleIdsRows, err := tx.Query(ctx, roleIdsQs.String(), roleIdsQs.Args()...)
			if err != nil {
				return &txRes, err
			}
			err = pgutils.ScanRowsToStructs(roleIdsRows, &roleIds)
			if err != nil {
				return &txRes, err
			}

			fmt.Printf(`roleIdsQs %s`, roleIdsQs.String())

			if len(roleIds) > 0 {
				if len(query.RoleIds) == 0 {
					query.RoleIds = make([]int64, 0)
				}
				for _, roleId := range roleIds {
					query.RoleIds = append(query.RoleIds, roleId.Id)
				}
			}
		}

		if len(query.RoleIds) > 0 {
			qs.AddWhere("AND", "role_ids", "&&", query.RoleIds)
		}

		if query.SearchTerm != "" {
			qs.StartGroupedWhere("AND").
				AddToGroupedWhere("OR", "first_name", "ILIKE", "%"+query.SearchTerm+"%", 0).
				AddToGroupedWhere("OR", "last_name", "ILIKE", nil, qs.CurrentWhereGroupIndex()).
				AddToGroupedWhere("OR", "phone_number", "ILIKE", nil, qs.CurrentWhereGroupIndex()).
				AddToGroupedWhere("OR", "email", "ILIKE", nil, qs.CurrentWhereGroupIndex()).
				CloseGroupedWhere()
		}

		// Get count query before applying pagination
		countQuery := qs.GetCountQuery()

		qs.Paginate(&query.PaginationQuery, false)

		fmt.Print(qs.String())

		elapsed := time.Since(timeSt)
		fmt.Println("GetUsersQueryHandler took %s nano seconds", elapsed.Nanoseconds())

		var appUsers []entity.AppUser
		usersRows, err := tx.Query(ctx, qs.String(), qs.Args()...)
		if err != nil {
			return &txRes, err
		}
		err = pgutils.ScanRowsToStructs(
			usersRows,
			&appUsers,
		)
		if err != nil {
			return &txRes, err
		}

		for i := range appUsers {
			appUsers[i].Password = ""
		}
		txRes.AppUsers = appUsers

		var totalCountStruct []struct {
			Count int64 `db:"count"`
		}

		countQueryRows, err := tx.Query(ctx, countQuery, qs.Args()...)
		if err != nil {
			return &txRes, err
		}
		err = pgutils.ScanRowsToStructs(
			countQueryRows,
			&totalCountStruct,
		)
		if err != nil {
			return &txRes, err
		}

		if len(totalCountStruct) != 0 {
			txRes.TotalCount = totalCountStruct[0].Count
		}

		return &txRes, nil
	})
	if err != nil {
		return result.Err(err)
	}

	result.Data = res

	return result.Ok()
}
