package lifetime

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ndodanli/backend-api/configs"
	uow "github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg/unit_of_work"
	"github.com/ndodanli/backend-api/pkg/infrastructure/services"
	"github.com/ndodanli/backend-api/pkg/logger"
)

var (
	LoggerSingleton      logger.ILogger
	UOWSingleton         uow.IUnitOfWork
	AppServicesSingleton *services.AppServices
)

func InitiateSingletons(appLogger logger.ILogger, db *pgxpool.Pool, cfg *configs.Config, redisService services.IRedisService) {
	LoggerSingleton = appLogger
	UOWSingleton = uow.NewUnitOfWork(db)
	AppServicesSingleton = InitializeAppServices(cfg, redisService)
}

func InitializeAppServices(cfg *configs.Config, redisService services.IRedisService) *services.AppServices {
	var appServices services.AppServices
	appServices.JWTService = services.NewJWTService(cfg.Auth)
	appServices.RedisService = redisService
	appServices.SendgridService = services.NewSendgridService(&cfg.Sendgrid)
	return &appServices
}
