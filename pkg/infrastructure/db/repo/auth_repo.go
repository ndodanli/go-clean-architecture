package repo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	httperr "github.com/ndodanli/backend-api/pkg/errors"
	"github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg"
	"go/types"
	"time"
)

type IAuthRepo interface {
	GetRefreshTokenWithUUID(tokenUUID uuid.UUID) (*RefreshTokenWithUUIDRepoRes, error)
	UpdateRefreshToken(tokenId int64, expiresAt time.Time, tokenUUID uuid.UUID) (*GetIdAndPasswordRepoRes, error)
	GetIdAndPasswordWithUsername(username string) (*GetIdAndPasswordRepoRes, error)
	UpsertRefreshToken(appUserId int64, expiresAt time.Time, refreshToken uuid.UUID) (*types.Nil, error)
}

type AuthRepo struct {
	db  *pgxpool.Pool
	ctx context.Context
	tm  *pg.TxSessionManager
}

func NewAuthRepo(db *pgxpool.Pool, ctx context.Context, tm *pg.TxSessionManager) IAuthRepo {
	return &AuthRepo{
		db:  db,
		ctx: ctx,
		tm:  tm,
	}
}

type RefreshTokenWithUUIDRepoRes struct {
	ID        int64     `json:"id"`
	AppUserId int64     `json:"appUserId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (r *AuthRepo) GetRefreshTokenWithUUID(tokenUUID uuid.UUID) (*RefreshTokenWithUUIDRepoRes, error) {
	return pg.ExecDefaultTx(r.ctx, r.tm, func(tx pgx.Tx) (*RefreshTokenWithUUIDRepoRes, error) {
		var res RefreshTokenWithUUIDRepoRes
		err := tx.QueryRow(r.ctx, `SELECT id, app_user_id, expires_at 
										FROM refresh_token 
										WHERE token_uuid = $1
										AND revoked = FALSE
										LIMIT 1`, tokenUUID).Scan(&res.ID, &res.AppUserId, &res.ExpiresAt)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, err
		}

		return &res, nil
	})
}

func (r *AuthRepo) UpdateRefreshToken(tokenId int64, expiresAt time.Time, tokenUUID uuid.UUID) (*GetIdAndPasswordRepoRes, error) {
	return pg.ExecDefaultTx(r.ctx, r.tm, func(tx pgx.Tx) (*GetIdAndPasswordRepoRes, error) {
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
		return nil, nil
	})
}

type GetIdAndPasswordRepoRes struct {
	ID       int64  `json:"id"`
	Password string `json:"password"`
}

func (r *AuthRepo) GetIdAndPasswordWithUsername(username string) (*GetIdAndPasswordRepoRes, error) {
	return pg.ExecDefaultTx(r.ctx, r.tm, func(tx pgx.Tx) (*GetIdAndPasswordRepoRes, error) {
		var res GetIdAndPasswordRepoRes
		err := tx.QueryRow(r.ctx, "SELECT id, password FROM app_user WHERE username = $1", username).Scan(&res.ID, &res.Password)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, err
		}

		return &res, nil
	})
}

func (r *AuthRepo) UpsertRefreshToken(appUserId int64, expiresAt time.Time, refreshToken uuid.UUID) (*types.Nil, error) {
	return pg.ExecDefaultTx(r.ctx, r.tm, func(tx pgx.Tx) (*types.Nil, error) {
		// Check if user's refresh token is revoked if it exists
		var revoked bool
		err := tx.QueryRow(r.ctx, `SELECT revoked FROM refresh_token WHERE app_user_id = $1`, appUserId).Scan(&revoked)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return nil, err
			}
		}

		if revoked {
			return nil, httperr.InvalidAuthenticationError
		}

		_, err = tx.Exec(r.ctx,
			`INSERT INTO refresh_token  (app_user_id, token_uuid, expires_at, created_at, updated_at)
    					VALUES ($1, $2, $3, NOW(), NOW())
    					ON CONFLICT (app_user_id) DO UPDATE SET
						    token_uuid = $2,
						    expires_at = $3,
						    updated_at = NOW()
						    `,
			appUserId, refreshToken, expiresAt)

		if err != nil {
			return nil, err
		}

		return nil, nil
	})
}
