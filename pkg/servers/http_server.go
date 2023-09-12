package servers

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/ndodanli/go-clean-architecture/api"
	"github.com/ndodanli/go-clean-architecture/configs"
	"github.com/ndodanli/go-clean-architecture/internal/auth"
	httpctrl "github.com/ndodanli/go-clean-architecture/internal/server/http/controller"
	res "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	cstmvalidator "github.com/ndodanli/go-clean-architecture/pkg/core/validator"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	jwtsvc "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/jwt"
	srvcns "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services/constants"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"github.com/swaggo/echo-swagger"
	"net/http"
	"strings"
	"time"
)

// @title Swagger Auth API
// @version 1.0
// @description This is an example server
// @contact.email ndodanli14@gmail.com
// @host 127.0.0.1:5005

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func (s *server) NewHttpServer(db *pgxpool.Pool, logger logger.Logger, auth *auth.Auth) (e *echo.Echo) {
	e = echo.New()

	// Handle ip extraction
	handleIpExtraction(e, s.cfg)

	// Timeout settings
	e.Use(middleware.TimeoutWithConfig(getTimeoutConfig()))

	// Request-Response middleware
	e.Use(getRequestResponseMiddleware(logger))

	// Global error handler
	e.HTTPErrorHandler = getGlobalErrorHandler(logger)

	// Recover from panics
	e.Use(middleware.RecoverWithConfig(getRecoverConfig()))

	// Gzip compression
	e.Use(middleware.GzipWithConfig(getGzipConfig()))

	// Decompress http requests if Content-Encoding header is set to gzip
	e.Use(middleware.DecompressWithConfig(getGzipDecompressConfig()))

	// CSRF protection
	e.Use(middleware.CSRF())

	// CQRS settings
	e.Use(middleware.CORSWithConfig(getCorsConfig()))

	// Set body limit
	e.Use(middleware.BodyLimit(bodyLimit))

	// Add request id to context
	e.Use(middleware.RequestID())

	// Request logger
	e.Use(middleware.RequestLoggerWithConfig(getLoggerConfig(logger)))

	// Security settings
	e.Use(middleware.SecureWithConfig(getSecureConfig()))

	// Validator settings
	e.Validator = cstmvalidator.NewCustomValidator(validator.New())

	// Swagger settings
	url := echoSwagger.URL("http://localhost:5005/swagger/doc.json")
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(url))

	//Versioning
	versionGroup := e.Group("/v1")

	// Authentication settings
	jwtService := jwtsvc.NewJwtService(s.cfg.Auth)
	versionGroup.Use(getJWTMiddleware(s.cfg, jwtService))

	// Register scoped instances(instances that are created per request)
	e.Use(registerScopedInstances(db))

	httpctrl.RegisterControllers(versionGroup, db)

	go func() {
		address := fmt.Sprintf("%s:%s", s.cfg.Http.HOST, s.cfg.Http.PORT)
		go func() {
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Routes:\n")
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

func getJWTMiddleware(cfg *configs.Config, jwtService jwtsvc.JwtServiceInterface) func(next echo.HandlerFunc) echo.HandlerFunc {
	validAudiences := strings.Split(cfg.Auth.JWT_AUDIENCES, ",")

	verifyAud := func(audiences string) bool {
		if validAudiences[0] == "*" {
			return true
		}
		for _, validAud := range validAudiences {
			for _, aud := range strings.Split(audiences, ",") {
				if validAud == aud {
					return true
				}
			}
		}
		return false
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := jwtService.ValidateToken(strings.Replace(c.Request().Header.Get("Authorization"), "Bearer ", "", 1))
			if err != nil {
				return httperr.UnauthorizedError
			}
			claims := token.Claims.(jwt.MapClaims)
			audiences := claims["aud"].(string)

			if !verifyAud(audiences) {
				return httperr.UnAuthorizedAudienceError
			}

			c.Set(srvcns.AuthUserKey, &jwtsvc.AuthUser{
				ID: claims["sub"].(string),
			})

			return next(c)
		}
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

func getGzipDecompressConfig() middleware.DecompressConfig {
	return middleware.DecompressConfig{
		Skipper: middleware.DefaultSkipper,
	}
}

func getTimeoutConfig() middleware.TimeoutConfig {
	return middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		Timeout:      timeout,
		ErrorMessage: "Request timed out. Please try again later.",
	}
}

func registerScopedInstances(db *pgxpool.Pool) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			txSessionManager := postgresql.NewTxSessionManager(db)
			c.Set(srvcns.TxSessionManagerKey, txSessionManager)
			return next(c)
		}
	}
}

func handleIpExtraction(e *echo.Echo, cfg *configs.Config) {
	switch strings.ToLower(cfg.Http.IP_EXTRACTION) {
	case "forwarded-for":
		e.IPExtractor = echo.ExtractIPFromXFFHeader()
	case "real-ip":
		e.IPExtractor = echo.ExtractIPFromRealIPHeader()
	case "no-proxy":
		e.IPExtractor = echo.ExtractIPDirect()
	default:
		e.IPExtractor = echo.ExtractIPDirect()
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
