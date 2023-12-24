package adminuserctrl

import (
	baseres "github.com/ndodanli/backend-api/pkg/core/response"
	"github.com/ndodanli/backend-api/pkg/infrastructure/mediatr"
	adminqueries "github.com/ndodanli/backend-api/pkg/infrastructure/mediatr/queries/admin"
)

func RegisterMediatrHandlers() error {
	var err error
	err = mediatr.RegisterRequestHandler[
		*adminqueries.GetUsersQuery, *baseres.Result[*adminqueries.GetUsersQueryResponse, error, struct{}],
	](&adminqueries.GetUsersQueryHandler{})
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[
		*adminqueries.UpdateUserRolesQuery, *baseres.Result[*adminqueries.UpdateUserRolesQueryResponse, error, struct{}],
	](&adminqueries.UpdateUserRolesQueryHandler{})
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[
		*adminqueries.BlockUsersQuery, *baseres.Result[*adminqueries.BlockUsersQueryResponse, error, struct{}],
	](&adminqueries.BlockUsersQueryHandler{})
	if err != nil {
		return err
	}

	return nil
}
