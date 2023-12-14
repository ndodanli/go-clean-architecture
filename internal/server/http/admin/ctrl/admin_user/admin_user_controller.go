package adminuserctrl

import (
	"github.com/labstack/echo/v4"
	_ "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/mediatr"
	adminqueries "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/mediatr/queries/admin"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"github.com/ndodanli/go-clean-architecture/pkg/utils"
	"net/http"
	"os"
)

type AdminUserController struct {
	cGroup     *echo.Group
	httpClient *http.Client
	logger     logger.ILogger
}

func NewAdminUserController(group *echo.Group, logger logger.ILogger) (*AdminUserController, error) {
	err := RegisterMediatrHandlers()
	if err != nil {
		return nil, err
	}
	if os.Getenv("APP_ENV") == "test" {
		return nil, err
	}
	c := &AdminUserController{
		cGroup: group.Group("/user"),
		logger: logger,
	}

	c.cGroup.GET("/getUsers", c.GetUsers)
	c.cGroup.POST("/updateRoles", c.UpdateUserRoles)

	return c, nil
}

func (ct *AdminUserController) GetUsers(c echo.Context) error {
	var query adminqueries.GetUsersQuery
	if err := utils.BindAndValidate(c, &query); err != nil {
		return err
	}
	res := mediatr.Send[*adminqueries.GetUsersQuery, *baseres.Result[*adminqueries.GetUsersQueryResponse, error, struct{}]](c, &query)
	if res.IsErr() {
		return res.GetErr()
	}
	return c.JSON(http.StatusOK, res)
}

func (ct *AdminUserController) UpdateUserRoles(c echo.Context) error {
	var command adminqueries.UpdateUserRolesQuery
	if err := utils.BindAndValidate(c, &command); err != nil {
		return err
	}
	res := mediatr.Send[*adminqueries.UpdateUserRolesQuery, *baseres.Result[*adminqueries.UpdateUserRolesQueryResponse, error, struct{}]](c, &command)
	if res.IsErr() {
		return res.GetErr()
	}
	return c.JSON(http.StatusOK, res)
}
