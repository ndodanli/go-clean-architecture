package router

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ndodanli/go-clean-architecture/pkg/adapter/controller"
)

func NewRouter(e *echo.Echo, c controller.AppController) *echo.Echo {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/users", func(context echo.Context) error { return c.User.GetUsers(context) })
	e.POST("/users", func(context echo.Context) error { return c.User.CreateUser(context) })

	return e
}
