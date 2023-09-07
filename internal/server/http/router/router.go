package httprouter

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	httpctrl "github.com/ndodanli/go-clean-architecture/internal/server/http/controller"
)

func NewRouter(e *echo.Echo, db *pgxpool.Pool) {
	httpctrl.RegisterControllers(e, db)
}
