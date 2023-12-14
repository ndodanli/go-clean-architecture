package servers

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/ndodanli/go-clean-architecture/api"
	"github.com/ndodanli/go-clean-architecture/configs"
	"github.com/ndodanli/go-clean-architecture/pkg/constant"
	res "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	cstmvalidator "github.com/ndodanli/go-clean-architecture/pkg/core/validator"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/pg"
	mw "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/middleware"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"github.com/ndodanli/go-clean-architecture/pkg/servers/lifetime"
	"github.com/swaggo/echo-swagger"
	"net/http"
	"os"
	"strings"
)

// @title Swagger Auth API
// @version 1.0
// @description This is an example server
// @contact.email ndodanli14@gmail.com
// @host 127.0.0.1:5005

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func (s *server) NewHttpServer(ctx context.Context, db *pgxpool.Pool, logger logger.ILogger, redisService *services.RedisService) (e *echo.Echo) {
	e = echo.New()

	// Initialize other middlewares
	mw.Init(s.cfg, lifetime.AppServicesSingleton, db)

	// Handle ip extraction
	handleIpExtraction(e, s.cfg)

	// Timeout settings
	e.Use(middleware.TimeoutWithConfig(getTimeoutConfig()))

	// Request-Response middleware
	e.Use(getRequestResponseMiddleware(lifetime.LoggerSingleton))

	// Global error handler
	e.HTTPErrorHandler = getGlobalErrorHandler(lifetime.LoggerSingleton)

	// Recover from panics
	e.Use(middleware.RecoverWithConfig(getRecoverConfig()))

	// Gzip compression
	e.Use(middleware.GzipWithConfig(getGzipConfig()))

	// Decompress http requests if Content-Encoding header is set to gzip
	e.Use(middleware.DecompressWithConfig(getGzipDecompressConfig()))

	// CSRF protection
	//e.Use(middleware.CSRFWithConfig(getCsrfConfig()))

	// CQRS settings
	e.Use(middleware.CORSWithConfig(getCorsConfig()))

	// Set body limit
	e.Use(middleware.BodyLimit(bodyLimit))

	// Trace Id
	e.Use(mw.TraceID)

	// Request logger
	e.Use(middleware.RequestLoggerWithConfig(getLoggerConfig(lifetime.LoggerSingleton)))

	// Security settings
	e.Use(middleware.SecureWithConfig(getSecureConfig()))

	// Set custom validator
	e.Validator = cstmvalidator.NewCustomValidator(validator.New())

	// Set custom binder
	//e.Binder = cstmbinder.NewCustomBinder()

	// Swagger settings
	url := echoSwagger.URL("http://localhost:5005/swagger/doc.json")
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(url))

	//Versioning
	versionGroup := e.Group("/v1")

	// Register scoped instances(instances that are created per req)
	e.Use(registerScopedInstances(db))

	err := RegisterControllers(versionGroup, lifetime.LoggerSingleton)
	if err != nil {
		lifetime.LoggerSingleton.Error("Error while registering controllers", err, "app")
		os.Exit(1)
	}

	go func() {
		address := fmt.Sprintf("%s:%s", s.cfg.Http.HOST, s.cfg.Http.PORT)
		routes := e.Routes()
		err := syncRoutesToDB(db, routes)
		err = insertDefaultAdminAppUser(db, &s.cfg.DefaultAdminAppUser)
		if err != nil {
			lifetime.LoggerSingleton.Error("Error while upserting routes to db", err, "app")
			os.Exit(1)
		}
		printRoutes(routes)
		go func() {
			select {
			case done := <-ctx.Done():
				lifetime.LoggerSingleton.Info(fmt.Sprintf("Server is shutting down. Reason: %s", done), nil, "app")
				if err := e.Shutdown(ctx); err != nil {
					lifetime.LoggerSingleton.Error("Server shutdown error", err, "app")
				}
			}
		}()

		e.Logger.Fatal(e.Start(address))

	}()

	return
}

func getGlobalErrorHandler(logger logger.ILogger) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		var he *echo.HTTPError
		if errors.As(err, &he) {
			errorData, ok := he.Message.(*httperr.ErrorData)
			if ok {
				baseHttpApiResult := res.NewResult[any, *echo.HTTPError, any](nil)
				if errorData.Status == 0 {
					errorData.Status = http.StatusInternalServerError
				}
				if errorData.Status >= 500 {
					baseHttpApiResult.SetErrorMessage("Internal Server Error")
				} else {
					baseHttpApiResult.SetErrorMessage(errorData.Message)
				}

				if errorData.ShouldLogAsError || errorData.Status >= 500 {
					logger.Error(errorData.Message, errorData.Metadata, c.Get(constant.General.TraceIDKey).(string))
				}
				if errorData.ShouldLogAsInfo {
					logger.Info(errorData.Message, errorData.Metadata, c.Get(constant.General.TraceIDKey).(string))
				}

				jsonError := c.JSON(errorData.Status, baseHttpApiResult)
				if jsonError != nil {
					logger.Error(err.Error(), err, c.Get(constant.General.TraceIDKey).(string))
				}
			} else {
				jsonError := c.JSON(he.Code, he.Message)
				if jsonError != nil {
					logger.Error(err.Error(), err, c.Get(constant.General.TraceIDKey).(string))
				}
			}
		} else {
			logger.Error(err.Error(), err, c.Get(constant.General.TraceIDKey).(string))
			result := res.NewResult[any, *echo.HTTPError, any](nil)
			result.SetErrorMessage("Internal Server Error")
			jsonError := c.JSON(http.StatusInternalServerError, result)
			if jsonError != nil {
				logger.Error(err.Error(), err, c.Get(constant.General.TraceIDKey).(string))
			}
		}
	}
}

func getRequestResponseMiddleware(logger logger.ILogger) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//logger.Info("Incoming req")

			err := next(c)
			if err != nil { //exec main process
				c.Error(err)
			}
			// Release all tx sessions if there are any
			txSessions := c.Get(constant.General.TxSessionManagerKey)
			if txSessions != nil {
				panicErr := txSessions.(*pg.TxSessionManager).ReleaseAllTxSessions(c.Request().Context(), err)
				if panicErr != nil {
					logger.Error("Error while releasing tx sessions", panicErr, c.Get(constant.General.TraceIDKey).(string))
				}
			}

			return nil
		}
	}
}

func getRecoverConfig() middleware.RecoverConfig {
	return middleware.RecoverConfig{
		Skipper:             middleware.DefaultSkipper,
		StackSize:           stackSize,
		DisableStackAll:     false,
		DisablePrintStack:   false,
		DisableErrorHandler: false,
		LogLevel:            0,
		LogErrorFunc:        nil,
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

func getLoggerConfig(logger logger.ILogger) middleware.RequestLoggerConfig {
	return middleware.RequestLoggerConfig{
		Skipper:      middleware.DefaultSkipper,
		LogURI:       true,
		LogStatus:    true,
		LogMethod:    true,
		LogRequestID: true,
		LogLatency:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info(fmt.Sprintf("HttpReq-Log M:%s, URI:%s, S:%d, L:%s", v.Method, v.URI, v.Status, v.Latency),
				nil, c.Get(constant.General.TraceIDKey).(string))

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

func getCsrfConfig() middleware.CSRFConfig {
	return middleware.CSRFConfig{
		Skipper: middleware.DefaultSkipper,
	}
}

func registerScopedInstances(db *pgxpool.Pool) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Session manager
			txSessionManager := pg.NewTxSessionManager(db)
			c.Set(constant.General.TxSessionManagerKey, txSessionManager)
			c.Set(constant.General.DBKey, db)
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

func syncRoutesToDB(db *pgxpool.Pool, routes []*echo.Route) error {
	ctx := context.Background()

	insertQuery := "INSERT INTO endpoint (name, method) VALUES "
	values := []interface{}{}

	for i, route := range routes {
		insertQuery += "(" + fmt.Sprintf("$%d, $%d)", 2*i+1, 2*i+2)
		values = append(values, route.Path, route.Method)

		if i < len(routes)-1 {
			insertQuery += ", "
		}
	}

	insertQuery += " ON CONFLICT (name, method) DO NOTHING"

	_, err := db.Exec(ctx, insertQuery, values...)
	if err != nil {
		return err
	}

	updateQuery := `
        UPDATE endpoint
        SET deleted_at = NOW()
        WHERE (name, method) NOT IN (`

	for i, _ := range routes {
		updateQuery += fmt.Sprintf("($%d, $%d)", 2*i+1, 2*i+2)

		if i < len(routes)-1 {
			updateQuery += ", "
		}
	}

	updateQuery += ") AND deleted_at = '0001-01-01 00:00:00'"

	_, err = db.Exec(ctx, updateQuery, values...)
	if err != nil {
		return err
	}

	return nil
}

func insertDefaultAdminAppUser(db *pgxpool.Pool, defaultAdminAppUser *configs.DefaultAdminAppUser) error {
	ctx := context.Background()
	var allEndpointIds []int64
	allEndpointIdsRows, err := db.Query(ctx, "SELECT id FROM endpoint")
	defer allEndpointIdsRows.Close()
	if err != nil {
		return err
	}
	for allEndpointIdsRows.Next() {
		var id int64
		err = allEndpointIdsRows.Scan(&id)
		if err != nil {
			return err
		}
		allEndpointIds = append(allEndpointIds, id)
	}

	// Insert default admin role
	var roleId int64
	err = db.QueryRow(ctx, "INSERT INTO role (name, endpoint_ids) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET endpoint_ids = $2 WHERE role.name = $1 RETURNING id", defaultAdminAppUser.USERNAME, allEndpointIds).Scan(&roleId)

	query := `
		INSERT INTO app_user (username, password, email, first_name, last_name, phone_number, role_ids)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (username) DO UPDATE SET
			password = $2,
			email = $3,
			first_name = $4,
			last_name = $5,
			phone_number = $6,
			role_ids = $7
		WHERE app_user.username = $1`

	_, err = db.Exec(ctx, query, defaultAdminAppUser.USERNAME, defaultAdminAppUser.PASSWORD, defaultAdminAppUser.EMAIL, defaultAdminAppUser.FIRST_NAME, defaultAdminAppUser.LAST_NAME, defaultAdminAppUser.PHONE_NUMBER, []int64{roleId})
	if err != nil {
		return err
	}

	return nil
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
