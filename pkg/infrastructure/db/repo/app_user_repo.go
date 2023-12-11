package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/app_user"
	"github.com/ndodanli/go-clean-architecture/pkg/utils"
	"go/types"
	"strings"
)

type IAppUserRepo interface {
	PatchAppUser(appUserID int64, updateProps map[string]interface{}) (*types.Nil, error)
	FindOneById(id int64, include []string) (*app_user.AppUser, error)
	FindOneByEmail(email string, include []string) (*app_user.AppUser, error)
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

func (r *AppUserRepo) FindOneByEmail(email string, include []string) (*app_user.AppUser, error) {
	return postgresql.ExecDefaultTx(r.ctx, r.tm, func(tx pgx.Tx) (*app_user.AppUser, error) {
		var res app_user.AppUser
		var query string
		if len(include) > 0 {
			query = `SELECT ` + strings.Join(include, ", ") + ` FROM app_user WHERE email = $1`
		} else {
			query = `SELECT * FROM app_user WHERE email = $1`
		}
		query += ` AND deleted_at = '0001-01-01 00:00:00' LIMIT 1`

		err := utils.ScanRowToStruct(
			tx.QueryRow(r.ctx, query, email),
			&res,
			include,
		)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, err
		}

		return &res, nil
	})
}

func (r *AppUserRepo) FindOneById(id int64, include []string) (*app_user.AppUser, error) {
	return postgresql.ExecDefaultTx(r.ctx, r.tm, func(tx pgx.Tx) (*app_user.AppUser, error) {
		var res app_user.AppUser
		var query string
		if len(include) > 0 {
			query = `SELECT ` + strings.Join(include, ", ") + ` FROM app_user WHERE id = $1`
		} else {
			query = `SELECT * FROM app_user WHERE id = $1`
		}

		query += ` AND deleted_at = '0001-01-01 00:00:00' LIMIT 1`

		row := tx.QueryRow(r.ctx, query, id)
		err := utils.ScanRowToStruct(row, &res, include)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, err
		}

		return &res, nil
	})

}

func (r *AppUserRepo) PatchAppUser(appUserID int64, updateProps map[string]interface{}) (*types.Nil, error) {
	return postgresql.ExecDefaultTx(r.ctx, r.tm, func(tx pgx.Tx) (*types.Nil, error) {
		updateQuery := "UPDATE app_user SET"
		values := []interface{}{appUserID}

		i := 2
		for key, value := range updateProps {
			key = utils.ToSnakeCase(key)
			updateQuery += fmt.Sprintf(" %s = $%d,", key, i)
			values = append(values, value)
			i++
		}

		updateQuery = updateQuery[:len(updateQuery)-1]

		updateQuery += " WHERE id = $1 AND deleted_at = '0001-01-01 00:00:00'"

		_, err := tx.Exec(r.ctx, updateQuery, values...)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, httperr.AppUserNotFoundError
			}
			return nil, err
		}

		return nil, nil
	})
}
