package httprouter

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	httpctrl "github.com/ndodanli/go-clean-architecture/internal/server/http/controller"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	serviceconstants "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services/constants"
)

func NewRouter(e *echo.Echo, db *pgxpool.Pool) {

}
