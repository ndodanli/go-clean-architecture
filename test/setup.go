package test

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ndodanli/backend-api/configs"
	httperr "github.com/ndodanli/backend-api/pkg/errors"
	"github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg"
	"github.com/ndodanli/backend-api/pkg/infrastructure/services"
	"github.com/ndodanli/backend-api/pkg/logger"
	"github.com/ndodanli/backend-api/pkg/servers"
	"github.com/ndodanli/backend-api/pkg/servers/lifetime"
	"os"
)

type TestEnv struct {
	Cfg           *configs.Config
	Ctx           context.Context
	RedisClient   *services.RedisService
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
	conn := pg.InitPgxPool(cfg, appLogger)

	//postgresql.Migrate(ctx, conn, appLogger)

	// Initialize http errors
	httperr.Init()

	// Initialize redis
	redisService := services.NewRedisService(cfg.Redis, appLogger)
	err = redisService.Ping(ctx)
	if err != nil {
		appLogger.Error(err.Error(), err, "app")
		//cancel()
	}
	defer func(client *services.RedisService) {
		err = client.Close()
		if err != nil {
			appLogger.Error(err.Error(), err, "app")
			//cancel()
		}
	}(redisService)

	lifetime.InitiateSingletons(appLogger, conn, cfg, redisService)

	appServices := lifetime.InitializeAppServices(cfg, redisService)

	err = servers.RegisterControllers(nil, appLogger)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &TestEnv{
		Cfg:           cfg,
		Ctx:           ctx,
		DB:            conn,
		Log:           appLogger,
		AppServices:   appServices,
		CancelContext: cancel,
	}
}

var (
	testEnv     *TestEnv
	cfg         *configs.Config
	ctx         context.Context
	db          *pgxpool.Pool
	log         *logger.ApiLogger
	appServices *services.AppServices
	ts          *pg.TxSessionManager
)

func setupTest() func() {
	// Setup
	testEnv = SetupTestEnv()
	cfg = testEnv.Cfg
	ctx = testEnv.Ctx
	db = testEnv.DB
	log = testEnv.Log
	appServices = testEnv.AppServices
	ts = pg.NewTxSessionManager(db)

	// Disable logs
	return func() {
		// Tear down
		defer db.Close()
		defer testEnv.CancelContext()
		txErr := ts.ReleaseAllTxSessionsForTestEnv(ctx, nil)
		if txErr != nil {
			fmt.Println(txErr)
		}
	}
}
