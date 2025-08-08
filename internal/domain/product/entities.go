package product

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"not null"`
	Type        string         `json:"type" gorm:"not null;check:type IN ('Mármore', 'Granito', 'Serviço')"`
	Price       float64        `json:"price" gorm:"not null"`
	Unit        string         `json:"unit" gorm:"not null"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
} 