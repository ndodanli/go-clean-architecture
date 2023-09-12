package uow

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/app_user/repos/app_user"
)

type UnitOfWorkInterface interface {
	GetDB() *pgxpool.Pool
	AppUserRepo(ctx context.Context) appuserrepo.IAppUserRepo
}

type UnitOfWork struct {
	db *pgxpool.Pool
}

func NewUnitOfWork(db *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{
		db: db,
	}
}

func (uow UnitOfWork) GetDB() *pgxpool.Pool {
	return uow.db
}

func (uow UnitOfWork) AppUserRepo(ctx context.Context) appuserrepo.IAppUserRepo {
	return appuserrepo.NewAppUserRepo(uow.db, ctx)
}
