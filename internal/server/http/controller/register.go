package httpctrl

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
)

type AppController struct {
	AuthController AuthControllerInterface
	echo           *echo.Echo
}

func RegisterControllers(e *echo.Echo, db *pgxpool.Pool) {
	var appServices services.AppServices
	appServices.AuthService = services.NewAuthService(uow.NewUnitOfWork(db))
	NewAuthController(e, appServices)
	//return &AppController{
	//	echo:           e,
	//	AuthController: NewAuthController(e, appServices),
	//}
}
