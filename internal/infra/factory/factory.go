package factory

import (
	"erp-api/internal/infra/database"
	"fmt"
	"os"
	"strings"
)

// NewRepositoryFactory creates a repository factory based on the database type
// Currently supports PostgreSQL, but can be extended for other databases
func NewRepositoryFactory(db database.Database) (database.RepositoryFactory, error) {
	dialect := strings.ToLower(strings.TrimSpace(os.Getenv("DB_DIALECT")))

	switch dialect {
	case "", "postgres", "postgresql":
		return NewPostgreSQLFactory(db)
	case "mysql":
		return NewMySQLFactory(db)
	default:
		return nil, fmt.Errorf("unsupported DB_DIALECT: %s", dialect)
	}
}
