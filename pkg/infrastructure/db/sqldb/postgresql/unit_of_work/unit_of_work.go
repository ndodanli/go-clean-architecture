package uow

import (
	"github.com/jackc/pgx/v4/pgxpool"
	repoports "github.com/ndodanli/go-clean-architecture/pkg/domain/ports/repositories"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/app_user/repos"
)

type UnitOfWorkInterface interface {
	GetDB() *pgxpool.Pool
	UserRepo() repoports.AppUserRepoInterface
	AcquireAllTxSessions() *postgresql.TxSessionManager
}

type UnitOfWork struct {
	db               *pgxpool.Pool
	txSessionManager *postgresql.TxSessionManager
}

func NewUnitOfWork(db *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{
		db: db,
	}
}

func (uow UnitOfWork) AcquireAllTxSessions() *postgresql.TxSessionManager {
	return uow.txSessionManager
}

func (uow UnitOfWork) GetDB() *pgxpool.Pool {
	return uow.db
}

func (uow UnitOfWork) UserRepo() repoports.AppUserRepoInterface {
	return repos.NewAppUserRepo(uow.db, uow.txSessionManager)
}
