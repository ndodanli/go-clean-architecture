package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ndodanli/go-clean-architecture/configs"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"os"
)

func InitPgxPool(cfg *configs.Config, logger logger.Logger) *pgxpool.Pool {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s pool_max_conns=%d", cfg.Postgresql.HOST, cfg.Postgresql.PORT, cfg.Postgresql.USER, cfg.Postgresql.PASS, cfg.Postgresql.DEFAULT_DB, cfg.Postgresql.MAX_CONN)
	conn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		logger.Error("PostgreSQL connection failed")
		os.Exit(1)
	}

	logger.Info("PostgreSQL connection established")

	return conn
}
