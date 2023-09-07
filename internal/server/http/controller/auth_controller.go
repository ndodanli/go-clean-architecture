package httpctrl

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"net/http"
)

type AuthControllerInterface interface {
	GetUser(c echo.Context) error
	UpdateTestString(c echo.Context) error
	RandomString(c echo.Context) error
}

type AuthController struct {
	echo        *echo.Echo
	group       *echo.Group
	httpClient  *http.Client
	authService services.AuthServiceInterface
	testString  string
}

func NewAuthController(echo *echo.Echo, requiredServices services.AppServices) AuthControllerInterface {
	ac := &AuthController{
		echo:        echo,
		group:       echo.Group("/auth"),
		authService: requiredServices.AuthService,
	}

	ac.group.GET("/user", ac.GetUser)
	ac.group.GET("/update", ac.UpdateTestString)
	ac.group.GET("/random", ac.RandomString)

	return ac
}

func (ac *AuthController) GetUser(c echo.Context) error {
	r := ac.authService.GetUser(c.Request().Context())
	//acMemoryAddress := fmt.Sprintf("%p", &ac)
	ac.testString = ac.testString + " fsdfds" + ac.testString
	resString := fmt.Sprintf("service string: %s, controller string: %s", r, ac.testString)
	return c.String(200, resString)
}

func (ac *AuthController) UpdateTestString(c echo.Context) error {
	str := ac.authService.UpdateTestString()
	return c.String(200, str)
}

func (ac *AuthController) RandomString(c echo.Context) error {
	str, _ := services.GenerateRandomString(10)
	return c.String(200, str)
}
