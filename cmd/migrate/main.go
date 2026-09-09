package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/migration"
)

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_DSN is required")
		os.Exit(2)
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()
	if err := migration.Up(context.Background(), db); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
