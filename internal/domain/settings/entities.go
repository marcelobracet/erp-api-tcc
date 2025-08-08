package settings

import (
	"time"

	"gorm.io/gorm"
)

type Settings struct {
	ID                string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TradeName         string         `json:"trade_name" gorm:"not null"`
	LegalName         string         `json:"legal_name" gorm:"not null"`
	CNPJ              string         `json:"cnpj" gorm:"not null"`
	Phone             string         `json:"phone" gorm:"not null"`
	Email             string         `json:"email" gorm:"not null"`
	Street            string         `json:"street" gorm:"not null"`
	Number            string         `json:"number" gorm:"not null"`
	Complement        string         `json:"complement,omitempty"`
	Neighborhood      string         `json:"neighborhood" gorm:"not null"`
	City              string         `json:"city" gorm:"not null"`
	State             string         `json:"state" gorm:"not null"`
	ZipCode           string         `json:"zip_code" gorm:"not null"`
	PrimaryColor      string         `json:"primary_color" gorm:"default:'#1976d2'"`
	SecondaryColor    string         `json:"secondary_color" gorm:"default:'#9c27b0'"`
	LogoURL           string         `json:"logo_url,omitempty"`
	CreatedAt         time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
} 