package authctrl

import (
	"github.com/labstack/echo/v4"
	_ "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/req"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	srvcns "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services/constants"
	"github.com/ndodanli/go-clean-architecture/pkg/utils"
	"net/http"
)

type AuthController struct {
	controllerGroup *echo.Group
	httpClient      *http.Client
	authService     services.AuthServiceInterface
}

func NewAuthController(group *echo.Group, requiredServices *services.AppServices) *AuthController {
	ac := &AuthController{
		controllerGroup: group.Group("/auth"),
		authService:     requiredServices.AuthService,
	}

	ac.controllerGroup.POST("/login", ac.Login)

	return ac
}

// Login godoc
// @Security BearerAuth
// @Summary      Login
// @Description  Login
// @Tags         app_user
// @Accept       json
// @Produce      json
// @Param        loginReq body req.LoginRequest true "Username"
// @Success      200  {object}   baseres.SwaggerSuccessRes[res.LoginRes] "OK. On success."
// @Failure      400  {object}   baseres.SwaggerValidationErrRes "Bad Request. On any validation error."
// @Failure      401  {object}   baseres.SwaggerUnauthorizedErrRes "Unauthorized."
// @Failure      500  {object}   baseres.SwaggerInternalErrRes "Internal Server Error."
// @Router       /v1/auth/login [post]
func (ac *AuthController) Login(c echo.Context) error {
	c.Response().Header().Set("Test-Header-Controller", "Test-Value Controller")
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
