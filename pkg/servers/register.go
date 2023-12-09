package servers

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/ndodanli/go-clean-architecture/configs"
	authctrl "github.com/ndodanli/go-clean-architecture/internal/server/http/ctrl/auth"
	testctrl "github.com/ndodanli/go-clean-architecture/internal/server/http/ctrl/test"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services/redissrv"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
)

type AppManager struct {
	*authctrl.AuthController
	*testctrl.TestController
	*echo.Echo
}

func RegisterControllers(e *echo.Group, db *pgxpool.Pool, cfg *configs.Config, redisService redissrv.IRedisService, logger logger.ILogger) {
	_, err := authctrl.NewAuthController(e, logger)
	if err != nil {
		return
	}
	_, err = testctrl.NewTestController(e, logger)
}
