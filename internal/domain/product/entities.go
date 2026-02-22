package product

import (
	"time"

	"erp-api/internal/utils/dbtypes"

	"gorm.io/gorm"
)

type Product struct {
	ID          dbtypes.UUID   `json:"id" gorm:"primaryKey"`
	TenantID    dbtypes.UUID   `json:"tenant_id" gorm:"not null"`
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

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = dbtypes.NewUUID()
	}
	return nil
}
