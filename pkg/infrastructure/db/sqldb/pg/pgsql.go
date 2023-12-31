package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ndodanli/backend-api/configs"
	"github.com/ndodanli/backend-api/pkg/logger"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	"os"
	"time"
)

func InitPgxPool(cfg *configs.Config, logger logger.ILogger) *pgxpool.Pool {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", cfg.Postgresql.HOST, cfg.Postgresql.PORT, cfg.Postgresql.USER, cfg.Postgresql.PASS, cfg.Postgresql.DEFAULT_DB)

	pgxConfig, err := pgxpool.ParseConfig(connStr)
	fmt.Printf("pgxConfig: %v\n", pgxConfig.ConnString())
	if err != nil {
		logger.Error("PostgreSQL connection failed", err, "app")
		os.Exit(1)
	}
	pgxConfig.MinConns = int32(cfg.Postgresql.MIN_CONN)
	pgxConfig.MaxConns = int32(cfg.Postgresql.MAX_CONN)
	pgxConfig.MaxConnLifetime = time.Duration(cfg.Postgresql.MAX_CONN_LIFETIME) * time.Second
	pgxConfig.MaxConnIdleTime = time.Duration(cfg.Postgresql.MAX_CONN_IDLE_TIME) * time.Second

	// Register pgxUUID type
	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	var conn *pgxpool.Pool
	conn, err = pgxpool.NewWithConfig(context.TODO(), pgxConfig)
	if err != nil {
		logger.Error("PostgreSQL connection failed", err, "app")
		os.Exit(1)
	}

	// Check connection
	err = conn.Ping(context.Background())
	if err != nil {
		logger.Error("PostgreSQL connection failed", err, "app")
		os.Exit(1)
	}

	logger.Info("PostgreSQL connection established", nil, "app")

	return conn
}
