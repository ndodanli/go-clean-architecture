package adminqueries

import (
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/pg"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/pg/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"github.com/ndodanli/go-clean-architecture/pkg/utils/pgutils"
)

type GetRolesAndEndpointsQueryHandler struct {
	UOW         uow.IUnitOfWork
	AppServices *services.AppServices
	Logger      logger.ILogger
	TM          *pg.TxSessionManager
}

type GetRolesAndEndpointsQuery struct {
}

type Role struct {
	RoleId      int64   `json:"roleId"`
	RoleName    string  `json:"roleName"`
	EndpointIds []int64 `json:"endpointIds"`
}

type RoleEndpoint struct {
	EndpointId     int64  `json:"endpointId"`
	EndpointName   string `json:"endpointName"`
	EndpointMethod string `json:"endpointMethod"`
}

type RolesAndEndpoints struct {
	RoleId    int64          `json:"roleId"`
	RoleName  string         `json:"roleName"`
	Endpoints []RoleEndpoint `json:"endpoints"`
}

type GetRolesAndEndpointsQueryResponse struct {
	RolesWithEndpoints []RolesAndEndpoints `json:"rolesWithEndpoints"`
	AllEndpoints       []RoleEndpoint      `json:"allEndpoints"`
}

func (h *GetRolesAndEndpointsQueryHandler) Handle(echoCtx echo.Context, query *GetRolesAndEndpointsQuery) *baseres.Result[*GetRolesAndEndpointsQueryResponse, error, struct{}] {
	result := baseres.NewResult[*GetRolesAndEndpointsQueryResponse, error, struct{}](&GetRolesAndEndpointsQueryResponse{})
	ctx := echoCtx.Request().Context()
	db := h.UOW.DB()

	var allEndpoints []RoleEndpoint
	allEndpointsRows, err := db.Query(ctx, `SELECT id, name, method FROM endpoint`)
	if err != nil {
		return result.Err(err)
	}
	err = pgutils.ScanRowsToStructs(allEndpointsRows, &allEndpoints)
	if err != nil {
		return result.Err(err)
	}
	result.Data.AllEndpoints = allEndpoints

	var allRoles []Role
	rolesRows, err := db.Query(ctx, `SELECT id, name, endpoint_ids FROM role`)
	if err != nil {
		return result.Err(err)
	}
	err = pgutils.ScanRowsToStructs(rolesRows, &allRoles)
	if err != nil {
		return result.Err(err)
	}

	var rolesWithEndpoints []RolesAndEndpoints
	for _, role := range allRoles {
		var endpoints []RoleEndpoint
		for _, endpointId := range role.EndpointIds {
			for _, endpoint := range allEndpoints {
				if endpoint.EndpointId == endpointId {
					endpoints = append(endpoints, endpoint)
				}
			}
		}
		rolesWithEndpoints = append(rolesWithEndpoints, RolesAndEndpoints{
			RoleId:    role.RoleId,
			RoleName:  role.RoleName,
			Endpoints: endpoints,
		})
	}
	result.Data.RolesWithEndpoints = rolesWithEndpoints

	return result.Ok()
}
