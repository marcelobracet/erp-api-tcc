package settings

import (
	"time"

	"erp-api/internal/utils/dbtypes"
	"gorm.io/gorm"
)

type Settings struct {
	ID        dbtypes.UUID `json:"id" gorm:"primaryKey"`
	TenantID  dbtypes.UUID `json:"tenant_id" gorm:"not null"`
	Key       string    `json:"key" gorm:"not null"`
	Value     string    `json:"value,omitempty"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (s *Settings) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = dbtypes.NewUUID()
	}
	return nil
}