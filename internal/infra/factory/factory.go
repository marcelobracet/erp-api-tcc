package factory

import (
	"erp-api/internal/infra/database"
)

// NewRepositoryFactory creates a repository factory based on the database type
// Currently supports PostgreSQL, but can be extended for other databases
func NewRepositoryFactory(db database.Database) (database.RepositoryFactory, error) {
	// Check database type and return appropriate factory
	// For now, we only support PostgreSQL
	return NewPostgreSQLFactory(db)
}


