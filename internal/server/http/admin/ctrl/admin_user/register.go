package adminuserctrl

import (
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/mediatr"
	adminqueries "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/mediatr/queries/admin"
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

	return nil
}
