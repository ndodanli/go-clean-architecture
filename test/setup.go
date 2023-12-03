package test

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ndodanli/go-clean-architecture/configs"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services/redissrv"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"github.com/ndodanli/go-clean-architecture/pkg/servers"
	"os"
)

type TestEnv struct {
	Cfg           *configs.Config
	Ctx           context.Context
	RedisClient   *redissrv.RedisService
	DB            *pgxpool.Pool
	Log           *logger.ApiLogger
	AppServices   *services.AppServices
	CancelContext context.CancelFunc
}

func SetupTestEnv() *TestEnv {
	err := os.Setenv("APP_ENV", "test")
	if err != nil {
		fmt.Println(err)
	}
	cfg, errConfig := configs.ParseConfig()
	if errConfig != nil {
		fmt.Println(errConfig)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Info(fmt.Sprintf("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.APP_VERSION, cfg.Logger.LEVEL, cfg.Server.APP_ENV), nil, "app")
	ctx, cancel := context.WithCancel(context.Background())

	conn := postgresql.InitPgxPool(cfg, appLogger)

	//postgresql.Migrate(ctx, conn, appLogger)

	// Initialize http errors
	httperr.Init()

	// Initialize redis
	client := redissrv.NewRedisService(cfg.Redis)
	err = client.Ping(ctx)
	if err != nil {
		appLogger.Error(err.Error(), err, "app")
		cancel()
	}
	defer func(client *redissrv.RedisService) {
		err = client.Close()
		if err != nil {
			appLogger.Error(err.Error(), err, "app")
			cancel()
		}
	}(client)

	appServices := servers.InitializeAppServices(conn, cfg, client.Client)

	return &TestEnv{
		Cfg:           cfg,
		Ctx:           ctx,
		RedisClient:   client,
		DB:            conn,
		Log:           appLogger,
		AppServices:   appServices,
		CancelContext: cancel,
	}
}
