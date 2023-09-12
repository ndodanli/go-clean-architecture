package httpctrl

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/ndodanli/go-clean-architecture/configs"
	authctrl "github.com/ndodanli/go-clean-architecture/internal/server/http/controller/auth"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	jwtsvc "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/jwt"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
)

type AppController struct {
	AuthController *authctrl.AuthController
	echo           *echo.Echo
}

func RegisterControllers(e *echo.Group, db *pgxpool.Pool, cfg *configs.Config) {
	appServices := initializeAppServices(db, cfg)
	authctrl.NewAuthController(e, appServices)
}

func initializeAppServices(db *pgxpool.Pool, cfg *configs.Config) *services.AppServices {
	var appServices services.AppServices
	appServices.JWTService = jwtsvc.NewJWTService(cfg.Auth)
	appServices.AuthService = services.NewAuthService(uow.NewUnitOfWork(db), appServices.JWTService)
	return &appServices
}
