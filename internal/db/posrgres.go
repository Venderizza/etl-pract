package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var PgPool *pgxpool.Pool

func ConnectPostgres() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := "localhost"

	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", user, password, host, dbName)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Postgres connect error:", err)
	}
	PgPool = pool
}
