package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/ndodanli/backend-api/configs"
	httperr "github.com/ndodanli/backend-api/pkg/errors"
	"github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg"
	oauthcfg "github.com/ndodanli/backend-api/pkg/infrastructure/oauth_cfg"
	"github.com/ndodanli/backend-api/pkg/infrastructure/services"
	"github.com/ndodanli/backend-api/pkg/logger"
	"github.com/ndodanli/backend-api/pkg/servers"
	"github.com/ndodanli/backend-api/pkg/servers/lifetime"
	"github.com/ndodanli/backend-api/pkg/utils/gracefulexit"
	"log"
)

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

	conn := pg.InitPgxPool(cfg, appLogger)

	pg.Migrate(ctx, conn, appLogger)

	newServer := servers.NewServer(cfg, &ctx, appLogger)

	// Initialize http errors
	httperr.Init()

	// Initialize oauth2 configs
	oauthcfg.Init(&cfg.GoogleOauth2)

	// Initialize redis
	redisService := services.NewRedisService(cfg.Redis, appLogger)
	err := redisService.Ping(ctx)
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

	newServer.NewHttpServer(ctx, conn, appLogger, redisService)

	// Exit from application gracefully
	gracefulexit.TerminateApp(ctx)

	appLogger.Info("Server Exited Properly", nil, "app")
}
