package httpctrl

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	httpctrl "github.com/ndodanli/go-clean-architecture/internal/server/http/controller/auth"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
)

type AppController struct {
	AuthController *httpctrl.AuthController
	echo           *echo.Echo
}

func RegisterControllers(e *echo.Group, db *pgxpool.Pool) {
	appServices := initializeAppServices(db)
	httpctrl.NewAuthController(e, appServices)
}

func initializeAppServices(db *pgxpool.Pool) *services.AppServices {
	var appServices services.AppServices
	appServices.AuthService = services.NewAuthService(uow.NewUnitOfWork(db))
	return &appServices
}
