package httpctrl

import (
	"github.com/labstack/echo/v4"
	httpctrl "github.com/ndodanli/go-clean-architecture/internal/server/http/controller/auth/req"
	_ "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"net/http"
)

type AuthControllerInterface interface {
	GetUser(c echo.Context) error
	PostUser(c echo.Context) error
}

type AuthController struct {
	echo            *echo.Group
	controllerGroup *echo.Group
	httpClient      *http.Client
	authService     services.AuthServiceInterface
	testString      string
}

func NewAuthController(group *echo.Group, requiredServices *services.AppServices) AuthControllerInterface {
	ac := &AuthController{
		echo:            group,
		controllerGroup: group.Group("/auth"),
		authService:     requiredServices.AuthService,
	}

	ac.controllerGroup.GET("/user", ac.GetUser)
	ac.controllerGroup.POST("/user", ac.PostUser)

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
// @Failure      500  {object}   res.SwaggerInternalErrRes "Internal Server Error."
// @Router       /v1/auth/user [get]
func (ac *AuthController) GetUser(c echo.Context) error {
	var reqParams httpctrl.GetUserRequest
	if err := c.Bind(&reqParams); err != nil {
		return err
	}

	if err := c.Validate(&reqParams); err != nil {
		return err
	}

	return c.String(200, "")
}

// PostUser godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id  body  GetUserRequest  true  "Account ID"
// @Success      200  {object}   res.SwaggerSuccessRes[GetUserResponse] "OK. On success."
// @Failure      400  {object}   res.SwaggerValidationErrRes "Bad Request. On any validation error."
// @Failure      500  {object}   res.SwaggerInternalErrRes "Internal Server Error."
// @Router       /v1/auth/user [post]
func (ac *AuthController) PostUser(c echo.Context) error {
	var reqParams httpctrl.GetUserRequest
	if err := c.Bind(&reqParams); err != nil {
		return err
	}

	if err := c.Validate(&reqParams); err != nil {
		return err
	}

	return c.String(200, "")
}
