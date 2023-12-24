package adminauthctrl

import (
	"github.com/labstack/echo/v4"
	_ "github.com/ndodanli/backend-api/pkg/core/response"
	baseres "github.com/ndodanli/backend-api/pkg/core/response"
	"github.com/ndodanli/backend-api/pkg/infrastructure/mediatr"
	"github.com/ndodanli/backend-api/pkg/infrastructure/mediatr/queries"
	adminqueries "github.com/ndodanli/backend-api/pkg/infrastructure/mediatr/queries/admin"
	"github.com/ndodanli/backend-api/pkg/logger"
	"github.com/ndodanli/backend-api/pkg/utils"
	"net/http"
	"os"
)

type AdminAuthController struct {
	cGroup     *echo.Group
	httpClient *http.Client
	logger     logger.ILogger
}

func NewAdminAuthController(group *echo.Group, logger logger.ILogger) (*AdminAuthController, error) {
	err := RegisterMediatrHandlers()
	if err != nil {
		return nil, err
	}
	if os.Getenv("APP_ENV") == "test" {
		return nil, err
	}
	c := &AdminAuthController{
		cGroup: group.Group("/auth"),
		logger: logger,
	}

	c.cGroup.POST("/login", c.Login)
	c.cGroup.GET("/refreshToken/:refreshToken", c.RefreshToken)
	c.cGroup.GET("/getRolesAndEndpoints", c.GetRolesAndEndpoints)
	c.cGroup.POST("/addOrUpdateRole", c.AddOrUpdateRole)
	c.cGroup.DELETE("/deleteRole/:roleId", c.DeleteRole)

	return c, nil
}

// Login godoc
// @Summary      Login
// @Description  Login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        loginReq body queries.LoginQuery true "Username"
// @Success      200  {object}   baseres.SwaggerSuccessRes[queries.LoginQueryResponse] "OK. On success."
// @Failure      400  {object}   baseres.SwaggerValidationErrRes "Bad Request. On any validation error."
// @Failure      401  {object}   baseres.SwaggerUnauthorizedErrRes "Unauthorized."
// @Failure      500  {object}   baseres.SwaggerInternalErrRes "Internal Server Error."
// @Router       /v1/auth/login [post]
func (ct *AdminAuthController) Login(c echo.Context) error {
	var query queries.LoginQuery
	if err := utils.BindAndValidate(c, &query); err != nil {
		return err
	}
	res := mediatr.Send[*queries.LoginQuery, *baseres.Result[*queries.LoginQueryResponse, error, struct{}]](c, &query)
	if res.IsErr() {
		return res.GetErr()
	}
	return c.JSON(http.StatusOK, res)
}

// RefreshToken godoc
// @Security BearerAuth
// @Summary      RefreshToken
// @Description  RefreshToken
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        refreshToken path string true "Refresh Token"
// @Success      200  {object}   baseres.SwaggerSuccessRes[queries.RefreshTokenQueryResponse] "OK. On success."
// @Failure      400  {object}   baseres.SwaggerValidationErrRes "Bad Request. On any validation error."
// @Failure      401  {object}   baseres.SwaggerUnauthorizedErrRes "Unauthorized."
// @Failure      500  {object}   baseres.SwaggerInternalErrRes "Internal Server Error."
// @Router       /v1/auth/refreshToken [get]
func (ct *AdminAuthController) RefreshToken(c echo.Context) error {
	var query queries.RefreshTokenQuery
	if err := utils.BindAndValidate(c, &query); err != nil {
		return err
	}
	res := mediatr.Send[*queries.RefreshTokenQuery, *baseres.Result[*queries.RefreshTokenQueryResponse, error, struct{}]](c, &query)
	if res.IsErr() {
		return res.GetErr()
	}
	return c.JSON(http.StatusOK, res)
}

func (ct *AdminAuthController) GetRolesAndEndpoints(c echo.Context) error {
	var query adminqueries.GetRolesAndEndpointsQuery
	res := mediatr.Send[*adminqueries.GetRolesAndEndpointsQuery, *baseres.Result[*adminqueries.GetRolesAndEndpointsQueryResponse, error, struct{}]](c, &query)
	if res.IsErr() {
		return res.GetErr()
	}
	return c.JSON(http.StatusOK, res)
}

func (ct *AdminAuthController) AddOrUpdateRole(c echo.Context) error {
	var query adminqueries.AddOrUpdateRoleQuery
	if err := utils.BindAndValidate(c, &query); err != nil {
		return err
	}
	res := mediatr.Send[*adminqueries.AddOrUpdateRoleQuery, *baseres.Result[*adminqueries.AddOrUpdateRoleQueryResponse, error, struct{}]](c, &query)
	if res.IsErr() {
		return res.GetErr()
	}
	return c.JSON(http.StatusOK, res)
}

func (ct *AdminAuthController) DeleteRole(c echo.Context) error {
	var query adminqueries.DeleteRoleQuery
	if err := utils.BindAndValidate(c, &query); err != nil {
		return err
	}
	res := mediatr.Send[*adminqueries.DeleteRoleQuery, *baseres.Result[*adminqueries.DeleteRoleQueryResponse, error, struct{}]](c, &query)
	if res.IsErr() {
		return res.GetErr()
	}
	return c.JSON(http.StatusOK, res)
}
