package settings

import (
	"time"
)

type Settings struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID  string    `json:"tenant_id" gorm:"type:uuid;not null"`
	Key       string    `json:"key" gorm:"not null"`
	Value     string    `json:"value,omitempty"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}