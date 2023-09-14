package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ndodanli/go-clean-architecture/configs"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	"os"
)

func InitPgxPool(cfg *configs.Config, logger logger.Logger) *pgxpool.Pool {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s pool_max_conns=%d", cfg.Postgresql.HOST, cfg.Postgresql.PORT, cfg.Postgresql.USER, cfg.Postgresql.PASS, cfg.Postgresql.DEFAULT_DB, cfg.Postgresql.MAX_CONN)

	pgxConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		logger.Error("PostgreSQL connection failed")
		os.Exit(1)
	}

	// Register pgxUUID type
	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	conn, err := pgxpool.NewWithConfig(context.TODO(), pgxConfig)
	if err != nil {
		logger.Error("PostgreSQL connection failed")
		os.Exit(1)
	}

	logger.Info("PostgreSQL connection established")

	return conn
}
