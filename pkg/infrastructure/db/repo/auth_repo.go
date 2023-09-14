package repo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	"time"
)

type IAuthRepo interface {
	GetRefreshTokenWithUUID(tokenUUID uuid.UUID, ts *postgresql.TxSessionManager) (*RefreshTokenWithUUIDRepoRes, error, uuid.UUID)
	UpdateRefreshToken(tokenId int64, expiresAt time.Time, tokenUUID uuid.UUID, sid uuid.UUID, ts *postgresql.TxSessionManager) (any, error)
	GetIdAndPasswordWithUsername(username string, ts *postgresql.TxSessionManager) (*GetOnlyIdRepoRes, error)
	CreateNewRefreshToken(appUserId int64, expiresAt time.Time, refreshToken uuid.UUID, ts *postgresql.TxSessionManager) (*RefreshTokenRepoRes, error)
}

type AuthRepo struct {
	db  *pgxpool.Pool
	ctx context.Context
}

func NewAuthRepo(db *pgxpool.Pool, ctx context.Context) IAuthRepo {
	return &AuthRepo{
		db:  db,
		ctx: ctx,
	}
}

func (r *AuthRepo) GetRefreshTokenWithUUID(tokenUUID uuid.UUID, ts *postgresql.TxSessionManager) (*RefreshTokenWithUUIDRepoRes, error, uuid.UUID) {
	tx, sid := ts.AcquireTxSession(r.ctx, uuid.Nil)
	var res RefreshTokenWithUUIDRepoRes
	err := tx.QueryRow(r.ctx, `SELECT id, app_user_id, expires_at 
										FROM refresh_token 
										WHERE token_uuid = $1
										AND revoked = FALSE
										LIMIT 1`, tokenUUID).Scan(&res.ID, &res.AppUserId, &res.ExpiresAt)

	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return nil, nil, sid
		}
		return nil, err, sid
	}

	return &res, nil, sid
}

func (r *AuthRepo) UpdateRefreshToken(tokenId int64, expiresAt time.Time, tokenUUID uuid.UUID, sid uuid.UUID, ts *postgresql.TxSessionManager) (any, error) {
	tx, _ := ts.AcquireTxSession(r.ctx, sid)
	defer ts.ReleaseTxSession(sid, r.ctx)
	_, err := tx.Exec(r.ctx,
		`UPDATE refresh_token 
					SET token_uuid = $1,
					    expires_at = $2,
					    updated_at = NOW()
                    WHERE id = $3`,
		tokenUUID, expiresAt, tokenId)

	if err != nil {
		return nil, err
	}

	err = tx.Commit(r.ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *AuthRepo) GetIdAndPasswordWithUsername(username string, ts *postgresql.TxSessionManager) (*GetOnlyIdRepoRes, error) {
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

func (r *AuthRepo) CreateNewRefreshToken(appUserId int64, expiresAt time.Time, refreshToken uuid.UUID, ts *postgresql.TxSessionManager) (*RefreshTokenRepoRes, error) {
	return postgresql.ExecTx(r.ctx, ts, uuid.Nil, func(tx pgx.Tx) (*RefreshTokenRepoRes, error) {
		_, err := tx.Exec(r.ctx,
			`INSERT INTO refresh_token (app_user_id, token_uuid, expires_at) 
					VALUES ($1, $2, $3)`,
			appUserId, refreshToken, expiresAt)

		if err != nil {
			return nil, err
		}

		return nil, nil
	})
}
