package adminuserctrl

import (
	"github.com/labstack/echo/v4"
	_ "github.com/ndodanli/backend-api/pkg/core/response"
	baseres "github.com/ndodanli/backend-api/pkg/core/response"
	"github.com/ndodanli/backend-api/pkg/infrastructure/mediatr"
	adminqueries "github.com/ndodanli/backend-api/pkg/infrastructure/mediatr/queries/admin"
	mw "github.com/ndodanli/backend-api/pkg/infrastructure/middleware"
	"github.com/ndodanli/backend-api/pkg/logger"
	"github.com/ndodanli/backend-api/pkg/utils"
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
		cGroup: group.Group("/user", mw.Auth),
		logger: logger,
	}

	c.cGroup.GET("/getUsers", c.GetUsers)
	c.cGroup.POST("/updateUserRoles", c.UpdateUserRoles)
	c.cGroup.POST("/blockUsers", c.BlockUsers)

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
	var query adminqueries.UpdateUserRolesQuery
	if err := utils.BindAndValidate(c, &query); err != nil {
		return err
	}
	res := mediatr.Send[*adminqueries.UpdateUserRolesQuery, *baseres.Result[*adminqueries.UpdateUserRolesQueryResponse, error, struct{}]](c, &query)
	if res.IsErr() {
		return res.GetErr()
	}
	return c.JSON(http.StatusOK, res)
}

func (ct *AdminUserController) BlockUsers(c echo.Context) error {
	var query adminqueries.BlockUsersQuery
	if err := utils.BindAndValidate(c, &query); err != nil {
		return err
	}
	res := mediatr.Send[*adminqueries.BlockUsersQuery, *baseres.Result[*adminqueries.BlockUsersQueryResponse, error, struct{}]](c, &query)
	if res.IsErr() {
		return res.GetErr()
	}
	return c.JSON(http.StatusOK, res)
}
