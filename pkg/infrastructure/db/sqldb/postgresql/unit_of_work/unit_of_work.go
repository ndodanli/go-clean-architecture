package uow

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/repo"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
)

type IUnitOfWork interface {
	GetDB() *pgxpool.Pool
	AppUserRepo(ctx context.Context, tm *postgresql.TxSessionManager) repo.IAppUserRepo
	AuthRepo(ctx context.Context, tm *postgresql.TxSessionManager) repo.IAuthRepo
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

func (uow UnitOfWork) AppUserRepo(ctx context.Context, tm *postgresql.TxSessionManager) repo.IAppUserRepo {
	return repo.NewAppUserRepo(uow.db, ctx, tm)
}

func (uow UnitOfWork) AuthRepo(ctx context.Context, tm *postgresql.TxSessionManager) repo.IAuthRepo {
	return repo.NewAuthRepo(uow.db, ctx, tm)
}
