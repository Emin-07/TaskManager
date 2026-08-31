package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// migrationsDir locates the repo-root `migrations` directory by walking up from
// the current working directory (the package under test) until it finds go.mod.
func migrationsDir(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, "migrations")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not find repo root (go.mod) from %s", dir)
		}
		dir = parent
	}
}

// TestDSN builds a postgres DSN from env vars, defaulting to values that match
// the docker compose setup when run locally.
func TestDSN() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		name = "postgres"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, name)
}

// OpenTestDB connects to postgres for integration tests, skipping the test
// when no database is reachable (e.g. running `go test -short`).
func OpenTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, err := sqlx.Open("postgres", TestDSN())
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		t.Skipf("skipping integration test, db not reachable: %v", err)
	}

	closeOnFail := func() { _ = db.Close() }

	if err := goose.SetDialect("postgres"); err != nil {
		closeOnFail()
		t.Fatalf("failed to set goose dialect: %v", err)
	}
	if err := goose.Up(db.DB, migrationsDir(t)); err != nil {
		closeOnFail()
		t.Fatalf("failed to run migrations: %v", err)
	}

	t.Cleanup(func() { _ = db.Close() })
	return db
}
