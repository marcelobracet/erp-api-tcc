package dbtypes

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// UUID is a cross-dialect UUID representation.
//
// - Postgres: uuid
// - MySQL:    char(36)
//
// It is defined as a named string type so it JSON-marshals as a string and is easy
// to use across the codebase.
type UUID string

func NewUUID() UUID {
	return UUID(uuid.NewString())
}

func (u UUID) String() string {
	return string(u)
}

func (UUID) GormDataType() string {
	return "uuid"
}

func (UUID) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	if db == nil || db.Dialector == nil {
		return "uuid"
	}
	name := strings.ToLower(strings.TrimSpace(db.Dialector.Name()))
	switch name {
	case "mysql", "mariadb":
		return "char(36)"
	case "postgres", "postgresql":
		return "uuid"
	default: // postgres, sqlite, etc
		return "uuid"
	}
}
