package quote

import (
	"time"

	"gorm.io/gorm"
)

type QuoteStatus string

const (
	QuoteStatusPending   QuoteStatus = "Pendente"
	QuoteStatusApproved  QuoteStatus = "Aprovado"
	QuoteStatusRejected  QuoteStatus = "Rejeitado"
	QuoteStatusCancelled QuoteStatus = "Cancelado"
)

type Quote struct {
	ID          string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ClientID    string         `json:"client_id" gorm:"not null"`
	Client      *Client        `json:"client,omitempty" gorm:"foreignKey:ClientID"`
	TotalValue  float64        `json:"total_value" gorm:"not null"`
	Status      QuoteStatus    `json:"status" gorm:"not null;default:'Pendente'"`
	Date        time.Time      `json:"date" gorm:"not null"`
	ValidUntil  time.Time      `json:"valid_until" gorm:"not null"`
	Notes       string         `json:"notes,omitempty"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type QuoteItem struct {
	ID        string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	QuoteID   string         `json:"quote_id" gorm:"not null"`
	ProductID string         `json:"product_id" gorm:"not null"`
	Product   *Product       `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Quantity  float64        `json:"quantity" gorm:"not null"`
	UnitPrice float64        `json:"unit_price" gorm:"not null"`
	TotalPrice float64       `json:"total_price" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Structs auxiliares para relacionamentos
type Client struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Type string  `json:"type"`
	Price float64 `json:"price"`
	Unit string  `json:"unit"`
} 