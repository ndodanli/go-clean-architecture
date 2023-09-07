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
	AcquireAllTxSessions() *postgresql.TxSessions
}

type UnitOfWork struct {
	db         *pgxpool.Pool
	txSessions *postgresql.TxSessions
}

func NewUnitOfWork(db *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{
		db:         db,
		txSessions: postgresql.NewTxSessions(db),
	}
}

func (uow UnitOfWork) AcquireAllTxSessions() *postgresql.TxSessions {
	return uow.txSessions
}

func (uow UnitOfWork) GetDB() *pgxpool.Pool {
	return uow.db
}

func (uow UnitOfWork) UserRepo() repoports.AppUserRepoInterface {
	return repos.NewAppUserRepo(uow.db, uow.txSessions)
}
