package quote

import (
	"time"

	"erp-api/internal/utils/dbtypes"
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
	ID       dbtypes.UUID `json:"id" gorm:"primaryKey"`
	TenantID dbtypes.UUID `json:"tenant_id" gorm:"not null;index"`
	ClientID dbtypes.UUID `json:"client_id" gorm:"not null;index"`
	UserID   dbtypes.UUID `json:"user_id" gorm:"not null;index"`

	Subtotal   float64 `json:"subtotal"`
	Discount   float64 `json:"discount"`
	TotalValue float64 `json:"total_value"`

	Status QuoteStatus `json:"status"`
	Notes  string      `json:"notes,omitempty"`

	ApprovedAt *time.Time `json:"approved_at,omitempty"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (q *Quote) BeforeCreate(tx *gorm.DB) error {
	if q.ID == "" {
		q.ID = dbtypes.NewUUID()
	}
	return nil
}

type QuoteItem struct {
	ID        dbtypes.UUID `json:"id" gorm:"primaryKey"`
	TenantID  dbtypes.UUID `json:"tenant_id" gorm:"not null"`
	QuoteID   dbtypes.UUID `json:"quote_id" gorm:"not null"`
	ProductID dbtypes.UUID `json:"product_id" gorm:"not null"`

	// Medidas
	WidthCM   float64 `json:"width_cm,omitempty"`  // largura
	HeightCM  float64 `json:"height_cm,omitempty"` // altura
	Thickness float64 `json:"thickness,omitempty"` // espessura
	AreaM2    float64 `json:"area_m2,omitempty"`   // calculado

	// Preço
	UnitPrice float64 `json:"unit_price" gorm:"not null"` // preço por m² ou unitário
	Quantity  int     `json:"quantity" gorm:"default:1"`
	Total     float64 `json:"total"` // calculado

	// Extras
	EdgeType       string `json:"edge_type,omitempty"` // borda
	HasCutout      bool   `json:"has_cutout" gorm:"default:false"`
	ReferenceImage string `json:"reference_image,omitempty"`
	Notes          string `json:"notes,omitempty"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (qi *QuoteItem) BeforeCreate(tx *gorm.DB) error {
	if qi.ID == "" {
		qi.ID = dbtypes.NewUUID()
	}
	return nil
}
