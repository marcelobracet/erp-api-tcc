package tenant

import (
	"time"
)

type Tenant struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"not null"`
	Plan      string    `json:"plan" gorm:"default:'free'"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
