package tenant

import (
	"time"

	"gorm.io/gorm"
)

type TenantStatus string

const (
	TenantStatusPending   TenantStatus = "pending"   // criado, ainda não pago
	TenantStatusActive    TenantStatus = "active"    // pago e liberado
	TenantStatusSuspended TenantStatus = "suspended" // inadimplente / bloqueado
	TenantStatusCanceled  TenantStatus = "canceled"
)

type Tenant struct {
	ID string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`

	// Identidade da empresa
	CompanyName string `json:"company_name" gorm:"not null"`
	TradeName   string `json:"trade_name,omitempty"` // nome fantasia
	CNPJ        string `json:"cnpj,omitempty" gorm:"index"`

	// Tipo de negócio (controlado, não validado)
	BusinessType string `json:"business_type" gorm:"default:'marmoraria'"`

	// Contato principal
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`

	// Plano & cobrança
	Plan   string       `json:"plan" gorm:"default:'free'"`
	Status TenantStatus `json:"status" gorm:"default:'pending'"`

	StripeCustomerID     string     `json:"stripe_customer_id,omitempty"`
	StripeSubscriptionID string     `json:"stripe_subscription_id,omitempty"`
	TrialEndsAt          *time.Time `json:"trial_ends_at,omitempty"`

	// Flags operacionais
	IsActive bool `json:"is_active" gorm:"default:true"`

	// Auditoria
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
