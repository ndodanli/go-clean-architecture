package singleton

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ndodanli/go-clean-architecture/configs"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services/redissrv"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
)

var (
	LoggerSingleton      logger.ILogger
	UOWSingleton         uow.IUnitOfWork
	AppServicesSingleton *services.AppServices
)

func InitiateSingletons(appLogger logger.ILogger, db *pgxpool.Pool, cfg *configs.Config, redisService redissrv.IRedisService) {
	LoggerSingleton = appLogger
	UOWSingleton = uow.NewUnitOfWork(db)
	AppServicesSingleton = InitializeAppServices(cfg, redisService)
}

func InitializeAppServices(cfg *configs.Config, redisService redissrv.IRedisService) *services.AppServices {
	var appServices services.AppServices
	appServices.JWTService = services.NewJWTService(cfg.Auth)
	appServices.RedisService = redisService
	return &appServices
}
