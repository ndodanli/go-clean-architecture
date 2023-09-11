package servers

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ndodanli/go-clean-architecture/configs"
	_ "github.com/ndodanli/go-clean-architecture/docs"
	"github.com/ndodanli/go-clean-architecture/internal/auth"
	httpctrl "github.com/ndodanli/go-clean-architecture/internal/server/http/controller"
	res "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	cstmvalidator "github.com/ndodanli/go-clean-architecture/pkg/core/validator"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
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

	// Request-Response middleware
	e.Use(getRequestResponseMiddleware(logger))

	// Global error handler
	e.HTTPErrorHandler = getGlobalErrorHandler(logger)

	// Recover from panics
	e.Use(middleware.RecoverWithConfig(getRecoverConfig()))

	// Gzip compression
	e.Use(middleware.GzipWithConfig(getGzipConfig()))

	// CSRF protection
	e.Use(middleware.CSRF())

	// CQRS setup
	e.Use(middleware.CORSWithConfig(getCorsConfig()))

	// Set body limit
	e.Use(middleware.BodyLimit(bodyLimit))

	// Add request id to context
	e.Use(middleware.RequestID())

	// Logger setup
	e.Use(middleware.RequestLoggerWithConfig(getLoggerConfig(logger)))

	// Security setup
	e.Use(middleware.SecureWithConfig(getSecureConfig()))

	// Validator setup
	e.Validator = cstmvalidator.NewCustomValidator(validator.New())

	// Swagger setup
	url := echoSwagger.URL("http://localhost:5005/swagger/doc.json")
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(url))

	//Versioning
	versionGroup := e.Group("/v1")

	// Authentication setup
	versionGroup.Use(echojwt.WithConfig(getJWTConfig(s.cfg)))

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

func getGlobalErrorHandler(logger logger.Logger) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		var he *echo.HTTPError
		if errors.As(err, &he) {
			errorData, ok := he.Message.(*httperr.ErrorData)
			if ok {
				baseHttpApiResult := res.NewResult[any, *echo.HTTPError, any]()
				baseHttpApiResult.SetErrorMessage(errorData.Message)
				if errorData.ShouldLog {
					logger.Error(err)
				}
				jsonError := c.JSON(he.Code, baseHttpApiResult)
				if jsonError != nil {
					logger.Error(err)
				}
			} else {
				jsonError := c.JSON(he.Code, he.Message)
				if jsonError != nil {
					logger.Error(err)
				}
			}
		} else {
			logger.Error(err)
			result := res.NewResult[any, *echo.HTTPError, any]()
			result.SetErrorMessage("Internal Server Error")
			jsonError := c.JSON(http.StatusInternalServerError, result)
			if jsonError != nil {
				logger.Error(err)
			}
		}
	}
}

func getRequestResponseMiddleware(logger logger.Logger) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//logger.Info("Incoming request")

			if err := next(c); err != nil { //exec main process
				c.Error(err)
			}

			//logger.Info("Outgoing response")

			return nil
		}
	}
}

func getJWTConfig(cfg *configs.Config) echojwt.Config {
	return echojwt.Config{
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().URL.Path, "swagger") {
				return true
			}
			return false
		},
		SigningKey: []byte(cfg.Auth.JWT_SECRET),
		ErrorHandler: func(c echo.Context, err error) error {
			return httperr.UnauthorizedError
		},
	}
}

func getRecoverConfig() middleware.RecoverConfig {
	return middleware.RecoverConfig{
		Skipper:           middleware.DefaultSkipper,
		StackSize:         stackSize,
		DisableStackAll:   false,
		DisablePrintStack: false,
		LogLevel:          0,
		LogErrorFunc:      nil,
	}
}

func getGzipConfig() middleware.GzipConfig {
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

func getCorsConfig() middleware.CORSConfig {
	return middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}
}

func getLoggerConfig(logger logger.Logger) middleware.RequestLoggerConfig {
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

func getSecureConfig() middleware.SecureConfig {
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
