package authctrl

import (
	"github.com/labstack/echo/v4"
	_ "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/mediatr"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/mediatr/queries"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"github.com/ndodanli/go-clean-architecture/pkg/utils"
	"net/http"
)

type AuthController struct {
	cGroup     *echo.Group
	httpClient *http.Client
	logger     logger.ILogger
}

func NewAuthController(group *echo.Group, logger logger.ILogger) (*AuthController, error) {
	err := RegisterMediatrHandlers()
	if err != nil {
		return nil, err
	}
	ac := &AuthController{
		cGroup: group.Group("/auth"),
		logger: logger,
	}

	ac.cGroup.POST("/login", ac.Login)
	ac.cGroup.GET("/refreshToken/:refreshToken", ac.RefreshToken)
	ac.cGroup.GET("/forgotPassword/:email", ac.ForgotPassword)

	return ac, nil
}

// Login godoc
// @Security BearerAuth
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
func (ac *AuthController) Login(c echo.Context) error {
	var query queries.LoginQuery
	if err := utils.BindAndValidate(c, &query); err != nil {
		return err
	}
	res := mediatr.Send[*queries.LoginQuery, *baseres.Result[queries.LoginQueryResponse, error, struct{}]](c, &query)
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
func (ac *AuthController) RefreshToken(c echo.Context) error {
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

// ForgotPassword godoc
// @Security BearerAuth
// @Summary      ForgotPassword
// @Description  ForgotPassword
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        ForgotPassword path string
// @Success      200  {object}   baseres.SwaggerSuccessRes[queries.ForgotPasswordQueryResponse] "OK. On success."
// @Failure      400  {object}   baseres.SwaggerValidationErrRes "Bad Request. On any validation error."
// @Failure      401  {object}   baseres.SwaggerUnauthorizedErrRes "Unauthorized."
// @Failure      500  {object}   baseres.SwaggerInternalErrRes "Internal Server Error."
// @Router       /v1/auth/ForgotPassword [get]
func (ac *AuthController) ForgotPassword(c echo.Context) error {
	var query queries.SendConfirmationEmailForgotPasswordQuery
	if err := utils.BindAndValidate(c, &query); err != nil {
		return err
	}
	res := mediatr.Send[*queries.SendConfirmationEmailForgotPasswordQuery, *baseres.Result[*queries.SendConfirmationEmailForgotPasswordQueryResponse, error, struct{}]](c, &query)
	if res.IsErr() {
		return res.GetErr()
	}
	return c.JSON(http.StatusOK, res)
}
