package repos

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	res "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	sqlresponses "github.com/ndodanli/go-clean-architecture/pkg/domain/ports/sql_responses"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewAppUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (ur *UserRepo) GetOneWithId(ctx context.Context, id int64) (sqlresponses.GetUserWithId, error) {
	row := ur.db.QueryRow(ctx, "SELECT id, name FROM app_user WHERE id = $1", id)
	var user sqlresponses.GetUserWithId
	// check if user exists

	err := row.Scan(&user.Id, &user.Name)

	if errors.Is(err, pgx.ErrNoRows) {
		return user, nil
	}

	return user, err
}

func (ur *UserRepo) InsertOne(ctx context.Context, name string) (sqlresponses.GetUserWithId, error) {
	var resUser sqlresponses.GetUserWithId

	tx, err := ur.db.Begin(ctx)
	if err != nil {
		return resUser, err
	}
	defer tx.Rollback(ctx)

	var rows pgx.Rows
	rows, err = tx.Query(ctx, "INSERT INTO app_user (name) VALUES ($1) RETURNING id, name", name)

	if err != nil {
		return resUser, err
	}

	for rows.Next() {
		err = rows.Scan(&resUser.Id, &resUser.Name)
		if err != nil {
			return resUser, err
		}
	}
	if err != nil {
		return resUser, err
	}

	// Commit the transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return resUser, err
	}

	return resUser, nil
}

type TestTxResult struct {
	ResultArray []struct {
		id   int64
		name string
	}
}

func (ur *UserRepo) TestTx(ctx context.Context) *res.Result[TestTxResult, error, any] {

	//	data, _ := postgresql.ExecTx(ctx, ur.txSessions, uuid.Nil, func(tx pgx.Tx) *res.Result[TestTxResult, error, any] {
	//		result := res.NewResult[TestTxResult, error, any]()
	//		query := `
	//select * from test_function($1);
	//`
	//		rows, err := tx.Query(ctx, query, "1")
	//		//rows.Close()
	//
	//		if err != nil {
	//			return result.Err(err)
	//		}
	//
	//		for rows.Next() {
	//			var row struct {
	//				id   int64
	//				name string
	//			}
	//
	//			err = rows.Scan(&row.id, &row.name)
	//			if err != nil {
	//				fmt.Printf("error: %v\n", err)
	//			}
	//
	//			result.Data.ResultArray = append(result.Data.ResultArray, row)
	//		}
	//
	//		if rows.Err() != nil {
	//			fmt.Printf("error: %v\n", rows.Err())
	//		}
	//
	//		rows, err = tx.Query(ctx, "SELECT * FROM test_function2($1)", "2")
	//
	//		if err != nil {
	//			fmt.Printf("error: %v\n", err)
	//		}
	//		rows.Close()
	//		if rows.Err() != nil {
	//			fmt.Printf("error: %v\n", rows.Err())
	//		}
	//
	//		err = tx.Commit(ctx)
	//
	//		if err != nil {
	//			fmt.Printf("error: %v\n", err)
	//		}
	//
	//		return result.Ok()
	//	})

	//return data
	return nil
}
