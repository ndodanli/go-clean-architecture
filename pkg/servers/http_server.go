package servers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/ndodanli/go-clean-architecture/docs"
	"github.com/ndodanli/go-clean-architecture/internal/auth"
	httpctrl "github.com/ndodanli/go-clean-architecture/internal/server/http/controller"
	cstmvalidator "github.com/ndodanli/go-clean-architecture/pkg/core/validator"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	serviceconstants "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services/constants"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"github.com/swaggo/echo-swagger"
	"net/http"
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

func (s *server) NewHttpServer(db *pgxpool.Pool, logger logger.Logger, auth *auth.Auth) (e *echo.Echo) {
	e = echo.New()

	// Recover from panics
	e.Use(middleware.RecoverWithConfig(NewRecoverConfig()))

	// Gzip compression
	e.Use(middleware.GzipWithConfig(NewGzipConfig()))

	// CSRF protection
	e.Use(middleware.CSRF())

	// CQRS setup
	e.Use(middleware.CORSWithConfig(NewCorsConfig()))

	// Set body limit
	e.Use(middleware.BodyLimit(bodyLimit))

	// Add request id to context
	e.Use(middleware.RequestID())

	// Logger setup
	e.Use(middleware.RequestLoggerWithConfig(NewLoggerConfig(logger)))

	// Security setup
	e.Use(middleware.SecureWithConfig(NewSecureConfig()))

	// Validator setup
	e.Validator = cstmvalidator.NewCustomValidator(validator.New())

	// Swagger setup
	url := echoSwagger.URL("http://localhost:5005/swagger/doc.json")
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(url))

	//Versioning
	versionGroup := e.Group("/v1")

	//e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	//	return func(c echo.Context) error {
	//		sub := "user1" // the user that wants to access a resource.
	//		obj := "data1" // the resource that is going to be accessed.
	//		act := "read"  // the operation that the user performs on the resource.
	//		dom := "tenant1"
	//		ok, err := auth.Enforcer().Enforce(sub, obj, act, dom)
	//		if err != nil {
	//			return err
	//		}
	//		if !ok {
	//			return echo.ErrForbidden
	//		}
	//		return next(c)
	//	}
	//})

	// Authentication setup
	versionGroup.Use(echojwt.WithConfig(NewJWTConfig()))

	// Register scoped instances(instances that are created per request)
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

func NewJWTConfig() echojwt.Config {
	return echojwt.Config{
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().URL.Path, "swagger") {
				return true
			}
			return false
		},
		SigningKey: []byte("secret"),
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		},
	}
}

func NewRecoverConfig() middleware.RecoverConfig {
	return middleware.RecoverConfig{
		Skipper:           middleware.DefaultSkipper,
		StackSize:         stackSize,
		DisableStackAll:   false,
		DisablePrintStack: false,
		LogLevel:          0,
		LogErrorFunc:      nil,
	}
}

func NewGzipConfig() middleware.GzipConfig {
	return middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().URL.Path, "swagger") {
				return true
			}
			return false
		},
		Level:     gzipLevel,
		MinLength: 0,
	}
}

func NewCorsConfig() middleware.CORSConfig {
	return middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}
}

func NewLoggerConfig(logger logger.Logger) middleware.RequestLoggerConfig {
	return middleware.RequestLoggerConfig{
		Skipper:      middleware.DefaultSkipper,
		LogURI:       true,
		LogStatus:    true,
		LogMethod:    true,
		LogRequestID: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Infof("Req-Log ID:%s M:%s URI:%s S:%d", v.RequestID, v.Method, v.URI, v.Status)

			return nil
		},
	}
}

func NewSecureConfig() middleware.SecureConfig {
	return middleware.SecureConfig{
		Skipper:            middleware.DefaultSkipper,
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
	}
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
