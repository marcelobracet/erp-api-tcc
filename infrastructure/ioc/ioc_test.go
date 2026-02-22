package ioc

import (
	"os"
	"testing"
)

func TestNewDatabase_SelectsMySQL(t *testing.T) {
	db, err := NewDatabase(DBConfig{Dialect: "mysql", Host: "localhost", Port: "3306", User: "u", Password: "p", Name: "n"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if db == nil {
		t.Fatalf("expected db instance, got nil")
	}
}

func TestNewDatabase_SelectsPostgresByDefault(t *testing.T) {
	db, err := NewDatabase(DBConfig{Dialect: "", Host: "localhost", Port: "5432", User: "u", Password: "p", Name: "n", SSLMode: "disable"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if db == nil {
		t.Fatalf("expected db instance, got nil")
	}
}

func TestReadDBConfigFromEnv_DefaultPorts(t *testing.T) {
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_SSLMODE", "")

	t.Setenv("DB_DIALECT", "mysql")
	cfg := ReadDBConfigFromEnv()
	if cfg.Port != "3306" {
		t.Fatalf("expected mysql default port 3306, got %s", cfg.Port)
	}

	t.Setenv("DB_DIALECT", "postgres")
	cfg = ReadDBConfigFromEnv()
	if cfg.Port != "5432" {
		t.Fatalf("expected postgres default port 5432, got %s", cfg.Port)
	}

	// Ensure env doesn't leak
	_ = os.Unsetenv("DB_DIALECT")
}
