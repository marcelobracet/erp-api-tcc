package quote

import (
	"time"

	"gorm.io/gorm"
)

type QuoteStatus string

const (
	QuoteStatusPending   QuoteStatus = "pending"
	QuoteStatusApproved  QuoteStatus = "approved"
	QuoteStatusRejected  QuoteStatus = "rejected"
	QuoteStatusCancelled QuoteStatus = "cancelled"
)

type Quote struct {
	ID             string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID       string         `json:"tenant_id" gorm:"not null"`
	ClientID       string         `json:"client_id" gorm:"not null"`
	UserID         string         `json:"user_id" gorm:"not null"`
	TotalValue     float64        `json:"total_value" gorm:"default:0"`
	Discount       float64        `json:"discount" gorm:"default:0"`
	Status         QuoteStatus    `json:"status" gorm:"default:'pending'"`
	ConversionRate *float64       `json:"conversion_rate,omitempty"`
	Notes          string         `json:"notes,omitempty"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type QuoteItem struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID  string    `json:"tenant_id" gorm:"not null"`
	QuoteID   string    `json:"quote_id" gorm:"not null"`
	ProductID string    `json:"product_id" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"not null;default:1"`
	Price     float64   `json:"price" gorm:"not null"`
	Total     float64   `json:"total" gorm:"->"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}