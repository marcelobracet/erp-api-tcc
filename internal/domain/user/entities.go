package user

import (
	"time"

	"erp-api/internal/utils/dbtypes"
	"gorm.io/gorm"
)

type User struct {
	ID          dbtypes.UUID `json:"id" gorm:"primaryKey"`
	KeycloakID  dbtypes.UUID `json:"keycloak_id" gorm:"column:keycloak_id;uniqueIndex;not null"`
	TenantID    dbtypes.UUID `json:"tenant_id" gorm:"not null;index"`
	DisplayName string       `json:"display_name" gorm:"column:display_name;not null"`
	Email       *string      `json:"email,omitempty" gorm:"column:email"`
	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	// In a Keycloak-based setup, the token subject is the stable identity.
	// Default to using it as our primary key to avoid a mapping layer.
	if u.ID == "" {
		if u.KeycloakID != "" {
			u.ID = u.KeycloakID
		} else {
			u.ID = dbtypes.NewUUID()
		}
	}
	if u.KeycloakID == "" {
		u.KeycloakID = u.ID
	}
	return nil
}

type CreateUserRequest struct {
	TenantID    string  `json:"tenant_id" validate:"required"`
	KeycloakID  string  `json:"keycloak_id" validate:"required"`
	DisplayName string  `json:"display_name" validate:"required,min=2"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
}

type UpdateUserRequest struct {
	DisplayName *string `json:"display_name,omitempty" validate:"omitempty,min=2"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
}