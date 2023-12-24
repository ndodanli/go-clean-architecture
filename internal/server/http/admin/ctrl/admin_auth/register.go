package adminauthctrl

import (
	baseres "github.com/ndodanli/backend-api/pkg/core/response"
	"github.com/ndodanli/backend-api/pkg/infrastructure/mediatr"
	"github.com/ndodanli/backend-api/pkg/infrastructure/mediatr/queries"
	adminqueries "github.com/ndodanli/backend-api/pkg/infrastructure/mediatr/queries/admin"
)

func RegisterMediatrHandlers() error {
	var err error
	err = mediatr.RegisterRequestHandler[
		*queries.LoginQuery, *baseres.Result[*queries.LoginQueryResponse, error, struct{}],
	](&queries.LoginQueryHandler{})
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[
		*queries.RefreshTokenQuery, *baseres.Result[*queries.RefreshTokenQueryResponse, error, struct{}],
	](&queries.RefreshTokenQueryHandler{})
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[
		*adminqueries.GetRolesAndEndpointsQuery, *baseres.Result[*adminqueries.GetRolesAndEndpointsQueryResponse, error, struct{}],
	](&adminqueries.GetRolesAndEndpointsQueryHandler{})
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[
		*adminqueries.AddOrUpdateRoleQuery, *baseres.Result[*adminqueries.AddOrUpdateRoleQueryResponse, error, struct{}],
	](&adminqueries.AddOrUpdateRoleQueryHandler{})
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[
		*adminqueries.DeleteRoleQuery, *baseres.Result[*adminqueries.DeleteRoleQueryResponse, error, struct{}],
	](&adminqueries.DeleteRoleQueryHandler{})
	if err != nil {
		return err
	}

	return nil
}
