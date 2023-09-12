package httpctrl

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	authctrl "github.com/ndodanli/go-clean-architecture/internal/server/http/controller/auth"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
)

type AppController struct {
	AuthController *authctrl.AuthController
	echo           *echo.Echo
}

func RegisterControllers(e *echo.Group, db *pgxpool.Pool) {
	appServices := initializeAppServices(db)
	authctrl.NewAuthController(e, appServices)
}

func initializeAppServices(db *pgxpool.Pool) *services.AppServices {
	var appServices services.AppServices
	appServices.AuthService = services.NewAuthService(uow.NewUnitOfWork(db))
	return &appServices
}
