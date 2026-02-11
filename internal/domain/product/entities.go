package product

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID    string         `json:"tenant_id" gorm:"type:uuid;not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description,omitempty"`
	Price       float64        `json:"price" gorm:"not null"`
	PriceType   string         `json:"price_type" gorm:"default:'unit'"`
	Stock       int            `json:"stock" gorm:"default:0"`
	SKU         string         `json:"sku,omitempty"`
	Category    string         `json:"category,omitempty"`
	ImageURL    string         `json:"image_url,omitempty"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
