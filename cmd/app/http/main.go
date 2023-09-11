package main

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/ndodanli/go-clean-architecture/configs"
	"github.com/ndodanli/go-clean-architecture/internal/auth"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"github.com/ndodanli/go-clean-architecture/pkg/servers"
	"github.com/ndodanli/go-clean-architecture/pkg/utils/gracefulexit"
	"log"
)

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
	defer conn.Close()

	newServer := servers.NewServer(cfg, &ctx, appLogger)

	a, err := auth.NewAuth(cfg)

	if err != nil {
		log.Fatalf("error: auth: %s", err)
	}

	// Initialize http errors
	httperr.Init()

	newServer.NewHttpServer(conn, appLogger, a)

	// Exit from application gracefully
	gracefulexit.TerminateApp(ctx)

	appLogger.Info("Server Exited Properly")
}
