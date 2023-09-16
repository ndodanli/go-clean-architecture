package main

import (
	"context"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/ndodanli/go-clean-architecture/configs"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	redissrv "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/cache/redis"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"github.com/ndodanli/go-clean-architecture/pkg/servers"
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

	cfg, errConfig := configs.ParseConfig()
	if errConfig != nil {
		log.Fatal(errConfig)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.APP_VERSION, cfg.Logger.LEVEL, cfg.Server.APP_ENV)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn := postgresql.InitPgxPool(cfg, appLogger)

	postgresql.Migrate(ctx, conn, appLogger)

	newServer := servers.NewServer(cfg, &ctx, appLogger)

	// Initialize http errors
	httperr.Init()

	// Initialize redis
	client := redissrv.NewRedisService(cfg.Redis)
	err := client.Ping(ctx)
	if err != nil {
		appLogger.Error(err)
		gracefulexit.TerminateApp(ctx)
	}
	defer func(client *redissrv.RedisService) {
		err = client.Close()
		if err != nil {
			appLogger.Error(err)
		}
	}(client)

	type data struct {
		Name string
	}
	var str []string = []string{"t"}
	var d data
	d, err = redissrv.AcquireHash(ctx, client.Client, "test1", str, func() (data, error) {
		return data{"test"}, nil
	})
	if err != nil {
		appLogger.Error(err)
	}
	appLogger.Info(d)

	newServer.NewHttpServer(conn, appLogger, client)

	// Exit from application gracefully
	gracefulexit.TerminateApp(ctx)

	appLogger.Info("Server Exited Properly")
}
