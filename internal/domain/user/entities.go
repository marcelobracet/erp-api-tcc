package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID     string         `json:"tenant_id" gorm:"type:uuid;not null"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null"`
	Password     string         `json:"-" gorm:"column:password_hash;not null"`
	Name         string         `json:"name" gorm:"not null"`
	Role         string         `json:"role" gorm:"not null;default:'user'"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateUserRequest struct {
	TenantID string `json:"tenant_id" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required,min=2"`
	Role     string `json:"role" validate:"required,oneof=admin user manager"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=2"`
	Role     *string `json:"role,omitempty" validate:"omitempty,oneof=admin user manager"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
} 