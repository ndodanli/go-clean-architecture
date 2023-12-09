package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	"go/types"
)

type IAppUserRepo interface {
	PatchAppUser(appUserID int64, updateProps map[string]interface{}) (*types.Nil, error)
}

type AppUserRepo struct {
	db  *pgxpool.Pool
	ctx context.Context
	tm  *postgresql.TxSessionManager
}

func NewAppUserRepo(db *pgxpool.Pool, ctx context.Context, tm *postgresql.TxSessionManager) IAppUserRepo {
	return &AppUserRepo{
		db:  db,
		ctx: ctx,
		tm:  tm,
	}
}

func (r *AppUserRepo) PatchAppUser(appUserID int64, updateProps map[string]interface{}) (*types.Nil, error) {
	return postgresql.ExecDefaultTx(r.ctx, r.tm, func(tx pgx.Tx) (*types.Nil, error) {
		updateQuery := "UPDATE app_user SET"
		values := []interface{}{appUserID}

		i := 2
		for key, value := range updateProps {
			updateQuery += fmt.Sprintf(" %s = $%d,", key, i)
			values = append(values, value)
			i++
		}

		updateQuery = updateQuery[:len(updateQuery)-1]

		updateQuery += " WHERE id = $1"

		_, err := tx.Exec(r.ctx, updateQuery, values...)
		if err != nil {
			if errors.As(err, &pgx.ErrNoRows) {
				return nil, httperr.AppUserNotFoundError
			}
			return nil, err
		}

		return nil, nil
	})
}
