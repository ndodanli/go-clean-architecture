package ctrl

import (
	"github.com/labstack/echo/v4"
	srvcns "github.com/ndodanli/go-clean-architecture/pkg/core/constant"
	_ "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/req"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/utils"
	"net/http"
)

type AuthController struct {
	controllerGroup *echo.Group
	httpClient      *http.Client
	authService     services.IAuthService
}

func NewAuthController(group *echo.Group, requiredServices *services.AppServices) *AuthController {
	ac := &AuthController{
		controllerGroup: group.Group("/auth"),
		authService:     requiredServices.AuthService,
	}

	ac.controllerGroup.POST("/login", ac.Login)
	ac.controllerGroup.GET("/refreshToken/:refreshToken", ac.RefreshToken)

	return ac
}

// Login godoc
// @Security BearerAuth
// @Summary      Login
// @Description  Login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        loginReq body req.LoginRequest true "Username"
// @Success      200  {object}   baseres.SwaggerSuccessRes[res.LoginRes] "OK. On success."
// @Failure      400  {object}   baseres.SwaggerValidationErrRes "Bad Request. On any validation error."
// @Failure      401  {object}   baseres.SwaggerUnauthorizedErrRes "Unauthorized."
// @Failure      500  {object}   baseres.SwaggerInternalErrRes "Internal Server Error."
// @Router       /v1/auth/login [post]
func (ac *AuthController) Login(c echo.Context) error {
	var payload req.LoginRequest
	if err := utils.BindAndValidate(c, &payload); err != nil {
		return err
	}

	result := ac.authService.Login(c.Request().Context(), payload, c.Get(srvcns.TxSessionManagerKey).(*postgresql.TxSessionManager))
	if result.IsError() {
		return result.GetError()
	}

	return c.JSON(http.StatusOK, result)
}

// RefreshToken godoc
// @Security BearerAuth
// @Summary      RefreshToken
// @Description  RefreshToken
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        refreshToken path string true "Refresh Token"
// @Success      200  {object}   baseres.SwaggerSuccessRes[res.RefreshTokenRes] "OK. On success."
// @Failure      400  {object}   baseres.SwaggerValidationErrRes "Bad Request. On any validation error."
// @Failure      401  {object}   baseres.SwaggerUnauthorizedErrRes "Unauthorized."
// @Failure      500  {object}   baseres.SwaggerInternalErrRes "Internal Server Error."
// @Router       /v1/auth/refreshToken [get]
func (ac *AuthController) RefreshToken(c echo.Context) error {
	var payload req.RefreshTokenRequest
	if err := utils.BindAndValidate(c, &payload); err != nil {
		return err
	}

	result := ac.authService.RefreshToken(c.Request().Context(), payload, c.Get(srvcns.TxSessionManagerKey).(*postgresql.TxSessionManager))
	if result.IsError() {
		return result.GetError()
	}

	return c.JSON(http.StatusOK, result)
}
