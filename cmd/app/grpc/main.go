package main

import (
	"context"
	"github.com/ndodanli/go-clean-architecture/configs"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"github.com/ndodanli/go-clean-architecture/pkg/servers"
	"github.com/ndodanli/go-clean-architecture/pkg/utils/gracefulexit"
	"google.golang.org/grpc/reflection"
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
	//appLogger.Info("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.APP_VERSION, cfg.Logger.LEVEL, cfg.Server.APP_ENV)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//conn = postgresql.InitPgxPool()
	//defer conn.Close()

	servers := servers.NewServer(cfg, &ctx, appLogger)

	grpcServer, errGrpcServer := servers.NewGrpcServer()
	if errGrpcServer != nil {
		cancel()
		return
	}

	if cfg.Server.APP_ENV == "dev" {
		reflection.Register(grpcServer)
	}

	// Exit from application gracefully
	gracefulexit.TerminateApp(ctx)

	grpcServer.GracefulStop()
	//appLogger.Info("Server Exited Properly")
}
