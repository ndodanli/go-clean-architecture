package repoports

import (
	"context"
	res "github.com/ndodanli/go-clean-architecture/pkg/core/respose"
	sqlresponses "github.com/ndodanli/go-clean-architecture/pkg/domain/ports/sql_responses"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/app_user/repos"
)

type AppUserRepoInterface interface {
	GetOneWithId(ctx context.Context, id int64) (sqlresponses.GetUserWithId, error)
	InsertOne(ctx context.Context, name string) (sqlresponses.GetUserWithId, error)
	TestTx(ctx context.Context) *res.Result[repos.TestTxResult, error, any]
}
