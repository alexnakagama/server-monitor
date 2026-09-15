package db

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func newTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("failed to load .env")
	}

	testDatabaseURL := os.Getenv("TEST_DATABASE_URL")
	if testDatabaseURL == "" {
		t.Fatal("TEST_DATABASE_URL is not set")
	}

	pool, err := NewPostgres(context.Background(), testDatabaseURL)
	if err != nil {
		t.Skipf("database is not available: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func TestNewPostgres(t *testing.T) {
	t.Run("returns error for invalid connection string", func(t *testing.T) {
		_, err := NewPostgres(context.Background(), "not-a-valid-url")
		if err == nil {
			t.Fatal("expected error for invalid URL, got nil")
		}
	})

	t.Run("returns error when database is unreachable", func(t *testing.T) {
		_, err := NewPostgres(context.Background(), "postgres://user:pass@localhost:1/testdb?connect_timeout=1")
		if err == nil {
			t.Fatal("expected error for unreachable database, got nil")
		}
	})

	t.Run("returns usable pool for valid connection", func(t *testing.T) {
		pool := newTestPool(t)

		conn, err := pool.Acquire(context.Background())
		if err != nil {
			t.Fatalf("failed to acquire connection: %v", err)
		}
		conn.Release()
	})
}
