package client

import (
	"time"

	"gorm.io/gorm"
)

type Client struct {
	ID           string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID     string         `json:"tenant_id" gorm:"type:uuid;not null"`
	Name         string         `json:"name" gorm:"not null"`
	Email        string         `json:"email,omitempty"`
	Phone        string         `json:"phone,omitempty"`
	Document     string         `json:"document" gorm:"not null"`
	DocumentType string         `json:"document_type" gorm:"not null;check:document_type IN ('CPF', 'CNPJ')"`
	Address      string         `json:"address,omitempty"`
	City         string         `json:"city,omitempty"`
	State        string         `json:"state,omitempty"`
	ZipCode      string         `json:"zip_code,omitempty"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
} 