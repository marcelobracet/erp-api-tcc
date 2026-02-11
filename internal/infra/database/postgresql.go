package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgreSQLDatabase implements the Database interface for PostgreSQL using GORM
type PostgreSQLDatabase struct {
	dsn    string
	db     *gorm.DB
	sqlDB  *sql.DB
	config *gorm.Config
}

// PostgreSQLConfig holds configuration for PostgreSQL connection
type PostgreSQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgreSQLDatabase creates a new PostgreSQL database instance
func NewPostgreSQLDatabase(config PostgreSQLConfig) Database {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	return &PostgreSQLDatabase{
		dsn: dsn,
		config: &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	}
}

// Connect establishes a connection to the PostgreSQL database
func (p *PostgreSQLDatabase) Connect(ctx context.Context) error {
	db, err := gorm.Open(postgres.Open(p.dsn), p.config)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	p.db = db

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	p.sqlDB = sqlDB

	// Set default connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return nil
}

// Close closes the database connection
func (p *PostgreSQLDatabase) Close() error {
	if p.sqlDB != nil {
		return p.sqlDB.Close()
	}
	return nil
}

// Ping checks if the database connection is alive
func (p *PostgreSQLDatabase) Ping(ctx context.Context) error {
	if p.sqlDB == nil {
		return fmt.Errorf("database connection is not established")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return p.sqlDB.PingContext(ctx)
}

// SetMaxOpenConns sets the maximum number of open connections
func (p *PostgreSQLDatabase) SetMaxOpenConns(n int) {
	if p.sqlDB != nil {
		p.sqlDB.SetMaxOpenConns(n)
	}
}

// SetMaxIdleConns sets the maximum number of idle connections
func (p *PostgreSQLDatabase) SetMaxIdleConns(n int) {
	if p.sqlDB != nil {
		p.sqlDB.SetMaxIdleConns(n)
	}
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
func (p *PostgreSQLDatabase) SetConnMaxLifetime(d time.Duration) {
	if p.sqlDB != nil {
		p.sqlDB.SetConnMaxLifetime(d)
	}
}

// AutoMigrate runs database migrations
func (p *PostgreSQLDatabase) AutoMigrate(dst ...any) error {
	if p.db == nil {
		return fmt.Errorf("database connection is not established")
	}
	return p.db.AutoMigrate(dst...)
}

// Exec executes a raw SQL query
func (p *PostgreSQLDatabase) Exec(sql string, values ...interface{}) error {
	if p.db == nil {
		return fmt.Errorf("database connection is not established")
	}
	return p.db.Exec(sql, values...).Error
}

// GetDB returns the underlying GORM DB instance
func (p *PostgreSQLDatabase) GetDB() interface{} {
	return p.db
}

// GetSQLDB returns the underlying sql.DB instance
func (p *PostgreSQLDatabase) GetSQLDB() (*sql.DB, error) {
	if p.sqlDB == nil {
		return nil, fmt.Errorf("database connection is not established")
	}
	return p.sqlDB, nil
}


