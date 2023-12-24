package authctrl

import (
	"github.com/labstack/echo/v4"
	_ "github.com/ndodanli/backend-api/pkg/core/response"
	baseres "github.com/ndodanli/backend-api/pkg/core/response"
	"github.com/ndodanli/backend-api/pkg/infrastructure/mediatr"
	"github.com/ndodanli/backend-api/pkg/infrastructure/mediatr/queries"
	oauthcfg "github.com/ndodanli/backend-api/pkg/infrastructure/oauth_cfg"
	"github.com/ndodanli/backend-api/pkg/logger"
	"github.com/ndodanli/backend-api/pkg/utils"
	"net/http"
	"os"
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
	if os.Getenv("APP_ENV") == "test" {
		return nil, err
	}
	c := &AuthController{
		cGroup: group.Group("/auth"),
		logger: logger,
	}

	//c.cGroup.POST("/register", c.Register)
	c.cGroup.POST("/login", c.Login)
	// TODO: Implement social logins for mobile apps
	c.cGroup.GET("/loginWithGoogle", c.LoginWithGoogle)
	c.cGroup.GET("/loginWithGoogle/callback", c.LoginWithGoogleCallback)
	c.cGroup.GET("/refreshToken/:refreshToken", c.RefreshToken)
	c.cGroup.GET("/forgotPassword/:email", c.ForgotPassword)
	c.cGroup.POST("/confirmForgotPasswordCode", c.ConfirmForgotPasswordCode)
	c.cGroup.GET("/emailConfirmation", c.EmailConfirmationHTML)
	c.cGroup.GET("/emailConfirmationConfirm", c.EmailConfirmationConfirm)

	return c, nil
}

// // Register godoc
// // @Summary      Register
// // @Description  Register
// // @Tags         Auth
// // @Accept       json
// // @Produce      json
// // @Param        loginReq body queries.RegisterQuery true "Username"
// // @Success      200  {object}   baseres.SwaggerSuccessRes[queries.RegisterQueryResponse] "OK. On success."
// // @Failure      400  {object}   baseres.SwaggerValidationErrRes "Bad Request. On any validation error."
// // @Failure      401  {object}   baseres.SwaggerUnauthorizedErrRes "Unauthorized."
// // @Failure      500  {object}   baseres.SwaggerInternalErrRes "Internal Server Error."
// // @Router       /v1/auth/login [post]
func (ac *AuthController) Register(c echo.Context) error {
	var query queries.RegisterQuery
	if err := utils.BindAndValidate(c, &query); err != nil {
		return err
	}
	res := mediatr.Send[*queries.RegisterQuery, *baseres.Result[*queries.RegisterQueryResponse, error, struct{}]](c, &query)
	if res.IsErr() {
		return res.GetErr()
	}
	return c.JSON(http.StatusOK, res)
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
func (ct *AuthController) Login(c echo.Context) error {
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
func (ct *AuthController) RefreshToken(c echo.Context) error {
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
func (ct *AuthController) ForgotPassword(c echo.Context) error {
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

// ConfirmForgotPasswordCode godoc
// @Summary      ConfirmForgotPasswordCode
// @Description  ConfirmForgotPasswordCode
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        ConfirmForgotPasswordCode path string
// @Success      200  {object}   baseres.SwaggerSuccessRes[queries.ConfirmForgotPasswordCodeQueryResponse] "OK. On success."
// @Failure      400  {object}   baseres.SwaggerValidationErrRes "Bad Request. On any validation error."
// @Failure      401  {object}   baseres.SwaggerUnauthorizedErrRes "Unauthorized."
// @Failure      500  {object}   baseres.SwaggerInternalErrRes "Internal Server Error."
// @Router       /v1/auth/ConfirmForgotPasswordCode [get]
func (ct *AuthController) ConfirmForgotPasswordCode(c echo.Context) error {
	var query queries.ConfirmForgotPasswordCodeQuery
	if err := utils.BindAndValidate(c, &query); err != nil {
		return err
	}
	res := mediatr.Send[*queries.ConfirmForgotPasswordCodeQuery, *baseres.Result[*queries.ConfirmForgotPasswordCodeQueryResponse, error, struct{}]](c, &query)
	if res.IsErr() {
		return res.GetErr()
	}
	return c.JSON(http.StatusOK, res)
}

func (ct *AuthController) LoginWithGoogle(c echo.Context) error {
	url := oauthcfg.GoogleOauth2Config.AuthCodeURL("state")
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (ct *AuthController) LoginWithGoogleCallback(c echo.Context) error {
	var query queries.LoginWithGoogleQuery
	if err := utils.BindAndValidate(c, &query); err != nil {
		return err
	}
	res := mediatr.Send[*queries.LoginWithGoogleQuery, *baseres.Result[*queries.LoginWithGoogleQueryResponse, error, struct{}]](c, &query)
	if res.IsErr() {
		return res.GetErr()
	}
	return c.JSON(http.StatusOK, res)
}

func (ct *AuthController) EmailConfirmationHTML(c echo.Context) error {
	htmlContent, err := os.ReadFile("assets/html/email_confirmation.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, htmlContent)
}

func (ct *AuthController) EmailConfirmationConfirm(c echo.Context) error {
	queryParam := c.QueryParam("code")
	_ = queryParam
	var query queries.EmailConfirmationQuery
	if err := utils.BindAndValidate(c, &query); err != nil {
		return err
	}
	res := mediatr.Send[*queries.EmailConfirmationQuery, *baseres.Result[*queries.EmailConfirmationQueryResponse, error, struct{}]](c, &query)
	if res.IsErr() {
		return res.GetErr()
	}
	return c.JSON(http.StatusOK, res)
}
