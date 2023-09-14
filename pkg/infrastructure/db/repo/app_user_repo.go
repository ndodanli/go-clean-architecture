package repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IAppUserRepo interface {
}

type AppUserRepo struct {
	db  *pgxpool.Pool
	ctx context.Context
}

func NewAppUserRepo(db *pgxpool.Pool, ctx context.Context) *AppUserRepo {
	return &AppUserRepo{
		db:  db,
		ctx: ctx,
	}
}
