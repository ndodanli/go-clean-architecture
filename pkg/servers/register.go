package servers

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/ndodanli/go-clean-architecture/configs"
	"github.com/ndodanli/go-clean-architecture/internal/server/http/ctrl"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/mediatr"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services/redissrv"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
)

type AppController struct {
	*ctrl.AuthController
	*ctrl.TestController
	*echo.Echo
}

func RegisterControllers(e *echo.Group, db *pgxpool.Pool, cfg *configs.Config, redisService redissrv.IRedisService, logger logger.ILogger) {
	unitOfWork := uow.NewUnitOfWork(db)
	appServices := InitializeAppServices(unitOfWork, cfg, redisService)
	ctrl.NewAuthController(e, appServices, logger)
	ctrl.NewTestController(e, logger)

	err := mediatr.RegisterMediatrHandlers(db, appServices, unitOfWork, logger)
	if err != nil {
		panic(err)
	}
}

func InitializeAppServices(uow uow.IUnitOfWork, cfg *configs.Config, redisService redissrv.IRedisService) *services.AppServices {
	var appServices services.AppServices
	appServices.JWTService = services.NewJWTService(cfg.Auth)
	appServices.RedisService = redisService
	return &appServices
}
