package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/ndodanli/go-clean-architecture/configs"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services/redissrv"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"github.com/ndodanli/go-clean-architecture/pkg/servers"
	"github.com/ndodanli/go-clean-architecture/pkg/servers/lifetime"
	"github.com/ndodanli/go-clean-architecture/pkg/utils/gracefulexit"
	"log"
)

type TestSt struct {
	Name    string
	Address string
	Phone   string
	Age     int
	Boolean bool
	Uuid    uuid.UUID
}

func main() {
	log.Println("Starting api server")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, errConfig := configs.ParseConfig()
	if errConfig != nil {
		log.Fatal(errConfig)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Info(fmt.Sprintf("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.APP_VERSION, cfg.Logger.LEVEL, cfg.Server.APP_ENV), nil, "app")

	conn := postgresql.InitPgxPool(cfg, appLogger)

	postgresql.Migrate(ctx, conn, appLogger)

	newServer := servers.NewServer(cfg, &ctx, appLogger)

	// Initialize http errors
	httperr.Init()

	// Initialize redis
	redisService := redissrv.NewRedisService(cfg.Redis, appLogger)
	err := redisService.Ping(ctx)
	if err != nil {
		appLogger.Error(err.Error(), err, "app")
		//cancel()
	}
	defer func(client *redissrv.RedisService) {
		err = client.Close()
		if err != nil {
			appLogger.Error(err.Error(), err, "app")
			//cancel()
		}
	}(redisService)

	lifetime.InitiateSingletons(appLogger, conn, cfg, redisService)

	newServer.NewHttpServer(ctx, conn, appLogger, redisService)

	// Exit from application gracefully
	gracefulexit.TerminateApp(ctx)

	appLogger.Info("Server Exited Properly", nil, "app")
}
