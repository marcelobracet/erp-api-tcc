package migrate

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Migrator encapsulates dialect-specific migration logic.
//
// Goal: keep the container free of database-specific SQL.
type Migrator interface {
	Run(db *gorm.DB) error
}

func NewMigrator(dialect string) (Migrator, error) {
	d := strings.ToLower(strings.TrimSpace(dialect))
	switch d {
	case "", "postgres", "postgresql":
		return &PostgresMigrator{}, nil
	case "mysql":
		return &MySQLMigrator{}, nil
	default:
		return nil, fmt.Errorf("unsupported DB_DIALECT for migrations: %s", dialect)
	}
}
