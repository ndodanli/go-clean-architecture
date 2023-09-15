package test

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ndodanli/go-clean-architecture/configs"
	redissrv "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/cache/redis"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/req"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"testing"
)

var (
	testEnv     *TestEnv
	cfg         *configs.Config
	ctx         context.Context
	redisClient *redissrv.RedisService
	db          *pgxpool.Pool
	log         *logger.ApiLogger
	appServices *services.AppServices
	ts          *postgresql.TxSessionManager
)

func setupTest() func(*error) {
	// Setup
	testEnv = SetupTestEnv()
	cfg = testEnv.Cfg
	ctx = testEnv.Ctx
	redisClient = testEnv.RedisClient
	db = testEnv.DB
	log = testEnv.Log
	appServices = testEnv.AppServices
	ts = postgresql.NewTxSessionManager(db)
	return func(err *error) {
		// Tear down
		defer db.Close()
		defer testEnv.CancelContext()
		txErr := ts.ReleaseAllTxSessionsForTestEnv(ctx, *err)
		if txErr != nil {
			fmt.Println(txErr)
		}
	}
}

func TestLogin(t *testing.T) {
	var err error
	defer setupTest()(&err)
	t.Run("should return jwt and refresh token", func(t *testing.T) {
		tableTestParams := []struct {
			name    string
			payload *req.LoginRequest
			want    bool
		}{
			{"fail authenticate", &req.LoginRequest{Username: "test", Password: "test1234"}, false},
			{"success authenticate", &req.LoginRequest{Username: "test", Password: "test123"}, true},
		}

		for _, param := range tableTestParams {
			t.Run(param.name, func(t *testing.T) {
				res := appServices.AuthService.Login(ctx, *param.payload, ts)
				got := res.IsSuccess()
				if got != param.want {
					err = res.GetError()
					t.Errorf("got %v want %v", got, param.want)
				}
			})
		}
	})
}
