package servers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/ndodanli/go-clean-architecture/docs"
	httpctrl "github.com/ndodanli/go-clean-architecture/internal/server/http/controller"
	cstmvalidator "github.com/ndodanli/go-clean-architecture/pkg/core/validator"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	serviceconstants "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services/constants"
	"github.com/swaggo/echo-swagger"
	"strings"
	"time"
)

const (
	ctxTimeout = 3
)

// @title Swagger Auth API
// @version 1.0
// @description This is a server for authentication and authorization

// @contact.email ndodanli14@gmail.com

// @host 127.0.0.1:5005

func (s *server) NewHttpServer(db *pgxpool.Pool) (e *echo.Echo) {
	e = echo.New()

	// Security setup
	e.Use(middleware.Secure())

	// Validator setup
	e.Validator = cstmvalidator.NewCustomValidator(validator.New())

	// Swagger setup
	url := echoSwagger.URL("http://localhost:5005/swagger/doc.json")
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(url))

	versionGroup := e.Group("/v1")

	e.Use(registerScopedInstances(db))

	httpctrl.RegisterControllers(versionGroup, db)

	go func() {
		address := fmt.Sprintf("%s:%s", s.cfg.Http.HOST, s.cfg.Http.PORT)
		go func() {
			// Wait for 300 ms
			time.Sleep(100 * time.Millisecond)
			printRoutes(e.Routes())
		}()
		e.Logger.Fatal(e.Start(address))

	}()

	return
}

func registerScopedInstances(db *pgxpool.Pool) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			txSessionManager := postgresql.NewTxSessionManager(db)
			c.Set(serviceconstants.TxSessionManagerKey, txSessionManager)
			return next(c)
		}
	}
}

func printRoutes(routes []*echo.Route) {
	routeMap := make(map[string][]string)
	for _, r := range routes {
		routeMap[r.Path] = append(routeMap[r.Path], r.Method)
	}
	colorGreen := "\033[0;32m"
	for path, methods := range routeMap {
		fmt.Printf("%s %s [%s] \n", colorGreen, path, strings.Join(methods, ", "))
	}
}
