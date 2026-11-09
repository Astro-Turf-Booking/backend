package database

import (
	"testing"
)

func TestConnect_MissingDatabaseURL(t *testing.T) {
	// Temporarily unset DATABASE_URL and clear env loading
	t.Setenv("DATABASE_URL", "")

	db, err := Connect()
	if err == nil {
		t.Fatalf("expected error when DATABASE_URL is empty, got nil db=%v", db)
	}

	if err.Error() != "DATABASE_URL environment variable is missing or empty" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestConnect_InvalidDatabaseURL(t *testing.T) {
	// Provide an invalid / unreachable connection URL
	t.Setenv("DATABASE_URL", "postgresql://invalid_user:invalid_pass@127.0.0.1:54329/invalid_db?sslmode=disable&connect_timeout=1")

	db, err := Connect()
	if err == nil {
		t.Fatalf("expected error when connecting to invalid database, got nil db=%v", db)
	}
}
