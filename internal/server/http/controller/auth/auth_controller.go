package httpctrl

import (
	"github.com/labstack/echo/v4"
	httpctrl "github.com/ndodanli/go-clean-architecture/internal/server/http/controller/auth/req"
	_ "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/utils"
	"net/http"
)

type AuthController struct {
	echo            *echo.Group
	controllerGroup *echo.Group
	httpClient      *http.Client
	authService     services.AuthServiceInterface
	testString      string
}

func NewAuthController(group *echo.Group, requiredServices *services.AppServices) *AuthController {
	ac := &AuthController{
		echo:            group,
		controllerGroup: group.Group("/auth"),
		authService:     requiredServices.AuthService,
	}

	ac.controllerGroup.GET("/user", ac.GetUser)

	return ac
}

// GetUser godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   query      int  true  "Account ID"
// @Success      200  {object}   res.SwaggerSuccessRes[GetUserResponse] "OK. On success."
// @Failure      400  {object}   res.SwaggerValidationErrRes "Bad Request. On any validation error."
// @Failure      401  {object}   res.SwaggerUnauthorizedErrRes "Unauthorized."
// @Failure      500  {object}   res.SwaggerInternalErrRes "Internal Server Error."
// @Router       /v1/auth/user [get]
func (ac *AuthController) GetUser(c echo.Context) error {
	c.Response().Header().Set("Test-Header-Controller", "Test-Value Controller")
	var reqParams httpctrl.GetUserRequest
	if err := utils.BindAndValidate(c, &reqParams); err != nil {
		return err
	}

	result := ac.authService.GetUser(1)

	if result.IsError() {
		//return result.GetError()
	}
	return c.JSON(http.StatusOK, reqParams)
}
