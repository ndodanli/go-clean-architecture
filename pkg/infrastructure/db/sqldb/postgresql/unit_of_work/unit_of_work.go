package uow

import (
	"github.com/jackc/pgx/v4/pgxpool"
	repoports "github.com/ndodanli/go-clean-architecture/pkg/domain/ports/repositories"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/app_user/repos"
)

type UnitOfWorkInterface interface {
	GetDB() *pgxpool.Pool
	UserRepo() repoports.AppUserRepoInterface
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

func (uow UnitOfWork) UserRepo() repoports.AppUserRepoInterface {
	return repos.NewAppUserRepo(uow.db)
}
