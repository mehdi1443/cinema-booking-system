package migration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestFreshDatabaseMigrationIntegration(t *testing.T) {
	dsn := os.Getenv("MIGRATION_INTEGRATION_DSN")
	if dsn == "" {
		t.Skip("MIGRATION_INTEGRATION_DSN is not set")
	}
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("parse migration DSN: %v", err)
	}
	databaseName := fmt.Sprintf("cinema_migration_test_%d", time.Now().UnixNano())
	adminConfig := config.Copy()
	adminConfig.Database = "postgres"
	adminDB, err := sql.Open("pgx", adminConfig.ConnString())
	if err != nil {
		t.Fatal(err)
	}
	defer adminDB.Close()
	if _, err := adminDB.ExecContext(context.Background(), "CREATE DATABASE "+databaseName); err != nil {
		t.Fatalf("create isolated test database: %v", err)
	}
	t.Cleanup(func() {
		_, _ = adminDB.Exec("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname=$1", databaseName)
		_, _ = adminDB.Exec("DROP DATABASE IF EXISTS " + databaseName)
	})
	testConfig := config.Copy()
	testConfig.Database = databaseName
	testDB, err := sql.Open("pgx", testConfig.ConnString())
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := Up(ctx, testDB); err != nil {
		t.Fatalf("migrate fresh database: %v", err)
	}
	if err := Up(ctx, testDB); err != nil {
		t.Fatalf("idempotent migration rerun: %v", err)
	}
	var versions, tables int
	if err := testDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM app_schema_migrations").Scan(&versions); err != nil {
		t.Fatal(err)
	}
	if err := testDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public' AND table_name IN ('users','payments','tickets','outbox_events','audit_logs','cinema_settlements')").Scan(&tables); err != nil {
		t.Fatal(err)
	}
	if versions < 10 || tables != 6 {
		t.Fatalf("incomplete migration result: versions=%d tables=%d", versions, tables)
	}
}
