package ioc

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"erp-api/internal/infra/container"
	"erp-api/internal/infra/database"
	"erp-api/internal/infra/factory"
	"erp-api/internal/infra/migrate"

	"gorm.io/gorm"
)

type DBConfig struct {
	Dialect  string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Params   string
}

func ReadDBConfigFromEnv() DBConfig {
	dialect := strings.ToLower(strings.TrimSpace(os.Getenv("DB_DIALECT")))

	cfg := DBConfig{
		Dialect:  dialect,
		Host:     strings.TrimSpace(os.Getenv("DB_HOST")),
		Port:     strings.TrimSpace(os.Getenv("DB_PORT")),
		User:     strings.TrimSpace(os.Getenv("DB_USER")),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     strings.TrimSpace(os.Getenv("DB_NAME")),
		SSLMode:  strings.TrimSpace(os.Getenv("DB_SSLMODE")),
		Params:   strings.TrimSpace(os.Getenv("DB_PARAMS")),
	}

	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == "" {
		switch cfg.Dialect {
		case "mysql", "mariadb":
			cfg.Port = "3306"
		default:
			cfg.Port = "5432"
		}
	}
	if cfg.User == "" {
		cfg.User = "erp_user"
	}
	if cfg.Password == "" {
		cfg.Password = "erp_password"
	}
	if cfg.Name == "" {
		cfg.Name = "erp_db"
	}
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}

	return cfg
}

func NewDatabase(cfg DBConfig) (database.Database, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Dialect)) {
	case "mysql", "mariadb":
		return database.NewMySQLDatabase(database.MySQLConfig{
			Host:     cfg.Host,
			Port:     cfg.Port,
			User:     cfg.User,
			Password: cfg.Password,
			DBName:   cfg.Name,
			Params:   cfg.Params,
		}), nil
	case "", "postgres", "postgresql":
		return database.NewPostgreSQLDatabase(database.PostgreSQLConfig{
			Host:     cfg.Host,
			Port:     cfg.Port,
			User:     cfg.User,
			Password: cfg.Password,
			DBName:   cfg.Name,
			SSLMode:  cfg.SSLMode,
		}), nil
	default:
		return nil, fmt.Errorf("unsupported DB_DIALECT: %s", cfg.Dialect)
	}
}

func ConnectWithRetry(ctx context.Context, db database.Database, maxWait time.Duration) error {
	deadline := time.Now().Add(maxWait)
	var lastErr error

	for attempt := 0; time.Now().Before(deadline); attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := db.Connect(attemptCtx)
		if err == nil {
			err = db.Ping(attemptCtx)
		}
		cancel()

		if err == nil {
			return nil
		}
		_ = db.Close()
		lastErr = err

		backoff := time.Duration(attempt+1) * 250 * time.Millisecond
		if backoff > 2*time.Second {
			backoff = 2 * time.Second
		}
		time.Sleep(backoff)
	}

	return fmt.Errorf("failed to connect to database after retries: %w", lastErr)
}

func BuildContainer(ctx context.Context) (*container.Container, error) {
	cfg := ReadDBConfigFromEnv()

	db, err := NewDatabase(cfg)
	if err != nil {
		return nil, err
	}

	if err := ConnectWithRetry(ctx, db, 45*time.Second); err != nil {
		return nil, err
	}

	gormDB, ok := db.GetDB().(*gorm.DB)
	if !ok {
		_ = db.Close()
		return nil, fmt.Errorf("failed to get GORM DB instance")
	}

	migrator, err := migrate.NewMigrator(cfg.Dialect)
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := migrator.Run(gormDB); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	repoFactory, err := factory.NewRepositoryFactory(db)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	app := container.NewContainer()
	app.Database = db
	app.RepoFactory = repoFactory
	app.DB = gormDB

	if err := app.InitializeWithPrewired(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return app, nil
}
