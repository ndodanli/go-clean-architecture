package ctrl

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/ndodanli/go-clean-architecture/configs"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/redis/go-redis/v9"
)

type AppController struct {
	AuthController *AuthController
	echo           *echo.Echo
}

func RegisterControllers(e *echo.Group, db *pgxpool.Pool, cfg *configs.Config, redisClient *redis.Client) {
	appServices := initializeAppServices(db, cfg, redisClient)
	NewAuthController(e, appServices)
}

func initializeAppServices(db *pgxpool.Pool, cfg *configs.Config, redisClient *redis.Client) *services.AppServices {
	var appServices services.AppServices
	appServices.JWTService = services.NewJWTService(cfg.Auth)
	appServices.AuthService = services.NewAuthService(uow.NewUnitOfWork(db), appServices.JWTService)
	return &appServices
}
