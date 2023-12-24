package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	httperr "github.com/ndodanli/backend-api/pkg/errors"
	"github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg"
	entity "github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg/entity/app_user"
	"github.com/ndodanli/backend-api/pkg/utils"
	"github.com/ndodanli/backend-api/pkg/utils/pgutils"
	"go/types"
	"strings"
)

type IAppUserRepo interface {
	PatchUser(appUserID int64, updateProps map[string]interface{}) (*types.Nil, error)
	FindOneById(id int64, include []string) (*entity.AppUser, error)
	FindManyByIds(ids []int64, include []string) ([]entity.AppUser, error)
	FindOneByEmail(email string, include []string) (*entity.AppUser, error)
	IsUserExist(id int64) (bool, error)
	BlockUsers(uid []int64, blockDetails entity.BlockDetails) ([]int64, error)
	IsUserBlocked(uid int64) (bool, error)
}

type AppUserRepo struct {
	db  *pgxpool.Pool
	ctx context.Context
	tm  *pg.TxSessionManager
}

func NewAppUserRepo(db *pgxpool.Pool, ctx context.Context, tm *pg.TxSessionManager) IAppUserRepo {
	return &AppUserRepo{
		db:  db,
		ctx: ctx,
		tm:  tm,
	}
}

func (r *AppUserRepo) IsUserBlocked(uid int64) (bool, error) {
	return pg.ExecDefaultTx(r.ctx, r.tm, func(tx pgx.Tx) (bool, error) {
		var res []struct {
			IsBlocked bool `db:"is_blocked"`
		}
		query := `WITH block_details_cte AS (
    SELECT
        app_user.id,
        elem->>'untilAt' AS until_at
    FROM
        app_user
            CROSS JOIN LATERAL unnest(app_user.block_details) AS elem
		WHERE id = $1
)
SELECT EXISTS (
    SELECT
        1
    FROM
        block_details_cte
    WHERE until_at::timestamptz > now()
) AS is_blocked
`

		rows, err := tx.Query(r.ctx, query, uid)
		if err != nil {
			return false, err
		}
		err = pgutils.ScanRowsToStructs(rows, &res)
		if err != nil {
			return false, err
		}

		if len(res) == 0 || !res[0].IsBlocked {
			return false, nil
		}

		return true, nil
	})
}

func (r *AppUserRepo) BlockUsers(uids []int64, blockDetails entity.BlockDetails) ([]int64, error) {
	return pg.ExecDefaultTx(r.ctx, r.tm, func(tx pgx.Tx) ([]int64, error) {
		var blockDetailsArr []entity.BlockDetails
		blockDetailsArr = append(blockDetailsArr, blockDetails)
		query := `UPDATE app_user SET block_details = block_details || $1 WHERE id = ANY($2) AND deleted_at = '0001-01-01 00:00:00' RETURNING id`
		rows, err := tx.Query(r.ctx, query, blockDetailsArr, uids)
		if err != nil {
			return nil, err
		}

		var res []pg.IdStruct
		err = pgutils.ScanRowsToStructs(rows, &res)
		if err != nil {
			return nil, err
		}

		return utils.ArrayMap(res, func(i pg.IdStruct) int64 {
			return i.Id
		}), nil
	})
}

func (r *AppUserRepo) IsUserExist(id int64) (bool, error) {
	return pg.ExecDefaultTx(r.ctx, r.tm, func(tx pgx.Tx) (bool, error) {
		var res []struct {
			Exists bool `db:"exists"`
		}
		query := `SELECT EXISTS(SELECT 1 FROM app_user WHERE id = $1 AND deleted_at = '0001-01-01 00:00:00')`

		rows, err := tx.Query(r.ctx, query, id)
		err = pgutils.ScanRowsToStructs(rows, &res)
		if err != nil {
			return false, err
		}

		if len(res) == 0 || !res[0].Exists {
			return false, nil
		}

		return true, nil
	})
}

func (r *AppUserRepo) FindOneByEmail(email string, include []string) (*entity.AppUser, error) {
	return pg.ExecDefaultTx(r.ctx, r.tm, func(tx pgx.Tx) (*entity.AppUser, error) {
		var res []entity.AppUser
		var query string
		if len(include) > 0 {
			query = `SELECT ` + strings.Join(include, ", ") + ` FROM app_user WHERE email = $1`
		} else {
			query = `SELECT * FROM app_user WHERE email = $1`
		}
		query += ` AND deleted_at = '0001-01-01 00:00:00' LIMIT 1`

		rows, err := tx.Query(r.ctx, query, email)
		err = pgutils.ScanRowsToStructs(
			rows,
			&res,
		)
		if err != nil {
			return nil, err
		}

		if len(res) == 0 {
			return nil, nil
		}

		return &res[0], nil
	})
}

func (r *AppUserRepo) FindManyByIds(ids []int64, include []string) ([]entity.AppUser, error) {
	return pg.ExecDefaultTx(r.ctx, r.tm, func(tx pgx.Tx) ([]entity.AppUser, error) {
		var res []entity.AppUser
		var query string
		if len(include) > 0 {
			query = `SELECT ` + strings.Join(include, ", ") + ` FROM app_user WHERE id = ANY($1)`
		} else {
			query = `SELECT * FROM app_user WHERE id = ANY($1)`
		}

		query += ` AND deleted_at = '0001-01-01 00:00:00'`

		rows, err := tx.Query(r.ctx, query, ids)
		err = pgutils.ScanRowsToStructs(
			rows,
			&res,
		)
		if err != nil {
			return nil, err
		}

		if len(res) == 0 {
			return nil, nil
		}

		return res, nil
	})
}

func (r *AppUserRepo) FindOneById(id int64, include []string) (*entity.AppUser, error) {
	return pg.ExecDefaultTx(r.ctx, r.tm, func(tx pgx.Tx) (*entity.AppUser, error) {
		var res []entity.AppUser
		var query string
		if len(include) > 0 {
			query = `SELECT ` + strings.Join(include, ", ") + ` FROM app_user WHERE id = $1`
		} else {
			query = `SELECT * FROM app_user WHERE id = $1`
		}

		query += ` AND deleted_at = '0001-01-01 00:00:00' LIMIT 1`

		rows, err := tx.Query(r.ctx, query, id)
		err = pgutils.ScanRowsToStructs(rows, &res)
		if err != nil {
			return nil, err
		}

		if len(res) == 0 {
			return nil, nil
		}

		return &res[0], nil
	})

}

func (r *AppUserRepo) PatchUser(appUserID int64, updateProps map[string]interface{}) (*types.Nil, error) {
	return pg.ExecDefaultTx(r.ctx, r.tm, func(tx pgx.Tx) (*types.Nil, error) {
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
