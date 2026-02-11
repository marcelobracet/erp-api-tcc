package database

import (
	"context"
	"database/sql"
	"time"

	clientDomain "erp-api/internal/domain/client"
	productDomain "erp-api/internal/domain/product"
	quoteDomain "erp-api/internal/domain/quote"
	settingsDomain "erp-api/internal/domain/settings"
	tenantDomain "erp-api/internal/domain/tenant"
	userDomain "erp-api/internal/domain/user"
)

type Database interface {
	// Connection management
	Connect(ctx context.Context) error
	Close() error
	Ping(ctx context.Context) error

	// Connection pool configuration
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	SetConnMaxLifetime(d time.Duration)

	// Migration support
	AutoMigrate(dst ...any) error
	Exec(sql string, values ...any) error

	// Get underlying database connection (for advanced operations)
	GetDB() any
	GetSQLDB() (*sql.DB, error)
}

// RepositoryFactory interface for creating repositories
// This allows different database implementations to provide their own repository factories
type RepositoryFactory interface {
	CreateTenantRepository() tenantDomain.Repository
	CreateUserRepository() userDomain.Repository
	CreateClientRepository() clientDomain.Repository
	CreateProductRepository() productDomain.Repository
	CreateQuoteRepository() quoteDomain.Repository
	CreateQuoteItemRepository() quoteDomain.ItemRepository
	CreateSettingsRepository() settingsDomain.Repository

	// Get the underlying database instance
	GetDatabase() Database
}
