package appuserrepo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
)

type IAppUserRepo interface {
	GetIdAndPasswordWithUsername(username string, ts *postgresql.TxSessionManager) (*GetOnlyIdRepoRes, error)
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

func (r *AppUserRepo) GetIdAndPasswordWithUsername(username string, ts *postgresql.TxSessionManager) (*GetOnlyIdRepoRes, error) {
	return postgresql.ExecTx(r.ctx, ts, uuid.Nil, func(tx pgx.Tx) (*GetOnlyIdRepoRes, error) {
		var res GetOnlyIdRepoRes
		err := tx.QueryRow(r.ctx, "SELECT id, password FROM app_user WHERE username = $1", username).Scan(&res.ID, &res.Password)

		if err != nil {
			if errors.As(err, &pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, err
		}

		return &res, nil
	})
}
