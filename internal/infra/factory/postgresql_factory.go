package factory

import (
	"fmt"

	clientDomain "erp-api/internal/domain/client"
	productDomain "erp-api/internal/domain/product"
	quoteDomain "erp-api/internal/domain/quote"
	settingsDomain "erp-api/internal/domain/settings"
	tenantDomain "erp-api/internal/domain/tenant"
	userDomain "erp-api/internal/domain/user"
	"erp-api/internal/infra/database"
	"erp-api/internal/infra/repository"

	"gorm.io/gorm"
)

// PostgreSQLFactory implements RepositoryFactory for PostgreSQL using GORM
type PostgreSQLFactory struct {
	db database.Database
}

// NewPostgreSQLFactory creates a new PostgreSQL repository factory
func NewPostgreSQLFactory(db database.Database) (database.RepositoryFactory, error) {
	if db == nil {
		return nil, fmt.Errorf("database cannot be nil")
	}

	return &PostgreSQLFactory{
		db: db,
	}, nil
}

// GetDatabase returns the underlying database instance
func (f *PostgreSQLFactory) GetDatabase() database.Database {
	return f.db
}

// getGormDB extracts the GORM DB instance from the Database interface
func (f *PostgreSQLFactory) getGormDB() (*gorm.DB, error) {
	dbInstance := f.db.GetDB()
	if dbInstance == nil {
		return nil, fmt.Errorf("database instance is nil")
	}

	gormDB, ok := dbInstance.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("database instance is not a GORM DB")
	}

	return gormDB, nil
}

// CreateTenantRepository creates a tenant repository
func (f *PostgreSQLFactory) CreateTenantRepository() tenantDomain.Repository {
	gormDB, err := f.getGormDB()
	if err != nil {
		// In a production system, you might want to handle this differently
		// For now, we'll panic as this indicates a configuration error
		panic(fmt.Sprintf("failed to get GORM DB: %v", err))
	}
	return repository.NewTenantRepository(gormDB)
}

// CreateUserRepository creates a user repository
func (f *PostgreSQLFactory) CreateUserRepository() userDomain.Repository {
	gormDB, err := f.getGormDB()
	if err != nil {
		panic(fmt.Sprintf("failed to get GORM DB: %v", err))
	}
	return repository.NewUserRepository(gormDB)
}

// CreateClientRepository creates a client repository
func (f *PostgreSQLFactory) CreateClientRepository() clientDomain.Repository {
	gormDB, err := f.getGormDB()
	if err != nil {
		panic(fmt.Sprintf("failed to get GORM DB: %v", err))
	}
	return repository.NewClientRepository(gormDB)
}

// CreateProductRepository creates a product repository
func (f *PostgreSQLFactory) CreateProductRepository() productDomain.Repository {
	gormDB, err := f.getGormDB()
	if err != nil {
		panic(fmt.Sprintf("failed to get GORM DB: %v", err))
	}
	return repository.NewProductRepository(gormDB)
}

// CreateQuoteRepository creates a quote repository
func (f *PostgreSQLFactory) CreateQuoteRepository() quoteDomain.Repository {
	gormDB, err := f.getGormDB()
	if err != nil {
		panic(fmt.Sprintf("failed to get GORM DB: %v", err))
	}
	return repository.NewQuoteRepository(gormDB)
}

// CreateQuoteItemRepository creates a quote item repository
func (f *PostgreSQLFactory) CreateQuoteItemRepository() quoteDomain.ItemRepository {
	gormDB, err := f.getGormDB()
	if err != nil {
		panic(fmt.Sprintf("failed to get GORM DB: %v", err))
	}
	return repository.NewQuoteItemRepository(gormDB)
}

// CreateSettingsRepository creates a settings repository
func (f *PostgreSQLFactory) CreateSettingsRepository() settingsDomain.Repository {
	gormDB, err := f.getGormDB()
	if err != nil {
		panic(fmt.Sprintf("failed to get GORM DB: %v", err))
	}
	return repository.NewSettingsRepository(gormDB)
}


