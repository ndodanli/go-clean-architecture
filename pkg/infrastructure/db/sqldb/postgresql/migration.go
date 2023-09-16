package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"os"
)

func Migrate(ctx context.Context, db *pgxpool.Pool, logger *logger.ApiLogger) {
	currentFilePath, _ := os.Getwd()
	migrationFilePath := currentFilePath + "/pkg/infrastructure/db/sqldb/postgresql/migration.sql"
	c, ioErr := os.ReadFile(migrationFilePath)
	if ioErr != nil {
		panic(ioErr)
	}
	sql := string(c)

	type TestError struct {
		Message string
		Test    string
		Boolean bool
	}

	t := TestError{
		Message: "test",
		Test:    "test",
		Boolean: true,
	}
	logger.Error("Test Log", t)
	_, err := db.Exec(ctx, sql)
	if err != nil {
		logger.Error("Test Log", t)
	}
}
