package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	httprouter "github.com/ndodanli/go-clean-architecture/internal/server/http/router"
	"os"
)

func main() {
	//ctx := context.Background()
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", "localhost", 5432, "postgres", "qweq", "postgres")
	conn, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	e := echo.New()

	httprouter.NewRouter(e, conn)

	e.Logger.Fatal(e.Start("127.0.0.1:1323"))

}
