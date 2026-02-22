package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Params   string // Example: parseTime=true&charset=utf8mb4&parseTime=True&loc=Local
}

type MySQLDatabase struct {
	dsn   string
	db    *gorm.DB
	sqlDB *sql.DB
	cfg   *gorm.Config
}

func NewMySQLDatabase(config MySQLConfig) Database {
	params := config.Params
	if params == "" {
		params = "parseTime=true&charset=utf8mb4&loc=Local"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		config.User, config.Password, config.Host, config.Port, config.DBName, params,
	)

	return &MySQLDatabase{
		dsn: dsn,
		cfg: &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	}
}

func (m *MySQLDatabase) Connect(ctx context.Context) error {
	db, err := gorm.Open(mysql.Open(m.dsn), m.cfg)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	m.db = db

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	m.sqlDB = sqlDB

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	return nil
}

func (m *MySQLDatabase) Close() error {
	if m.sqlDB != nil {
		return m.sqlDB.Close()
	}
	return nil
}

func (m *MySQLDatabase) Ping(ctx context.Context) error {
	if m.sqlDB == nil {
		return fmt.Errorf("database not connected")
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return m.sqlDB.PingContext(ctx)
}

func (m *MySQLDatabase) SetMaxOpenConns(n int) {
	if m.sqlDB != nil {
		m.sqlDB.SetMaxOpenConns(n)
	}
}

func (m *MySQLDatabase) SetMaxIdleConns(n int) {
	if m.sqlDB != nil {
		m.sqlDB.SetMaxIdleConns(n)
	}
}

func (m *MySQLDatabase) SetConnMaxLifetime(d time.Duration) {
	if m.sqlDB != nil {
		m.sqlDB.SetConnMaxLifetime(d)
	}
}

func (m *MySQLDatabase) AutoMigrate(dst ...any) error {
	return m.db.AutoMigrate(dst...)
}

func (m *MySQLDatabase) Exec(sql string, values ...any) error {
	return m.db.Exec(sql, values...).Error
}

func (m *MySQLDatabase) GetDB() any {
	return m.db
}

func (m *MySQLDatabase) GetSQLDB() (*sql.DB, error) {
	if m.sqlDB == nil {
		return nil, fmt.Errorf("database not connected")
	}
	return m.sqlDB, nil
}
