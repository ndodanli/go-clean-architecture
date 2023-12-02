package mediatr

import (
	"github.com/jackc/pgx/v5/pgxpool"
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/mediatr/queries"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
)

func RegisterMediatrHandlers(db *pgxpool.Pool, appServices *services.AppServices, uow uow.IUnitOfWork, logger logger.ILogger) error {
	var err error
	err = RegisterRequestHandler[
		*queries.LoginQuery, *baseres.Result[queries.LoginQueryResponse, error, struct{}],
	](queries.NewLoginQueryHandler(appServices, uow, logger))
	if err != nil {
		return err
	}

	err = RegisterRequestHandler[
		*queries.RefreshTokenQuery, *baseres.Result[queries.RefreshTokenQueryResponse, error, struct{}],
	](queries.NewRefreshTokenQueryHandler(appServices, uow, logger))
	if err != nil {
		return err
	}

	return nil
}
