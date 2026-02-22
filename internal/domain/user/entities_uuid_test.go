package user

import (
	"testing"

	"erp-api/internal/utils/dbtypes"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestUserBeforeCreate_SetsUUID(t *testing.T) {
	u := &User{}
	if err := u.BeforeCreate(&gorm.DB{}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if u.ID == "" {
		t.Fatalf("expected ID to be set")
	}
	if _, err := uuid.Parse(string(u.ID)); err != nil {
		t.Fatalf("expected valid uuid, got %q: %v", string(u.ID), err)
	}
	if u.KeycloakID == "" {
		t.Fatalf("expected KeycloakID to be set")
	}
}

func TestUserBeforeCreate_DoesNotOverrideExistingID(t *testing.T) {
	existing := uuid.NewString()
	u := &User{ID: dbtypes.UUID(existing)}
	_ = u.BeforeCreate(&gorm.DB{})
	if u.ID != dbtypes.UUID(existing) {
		t.Fatalf("expected ID unchanged")
	}
}
