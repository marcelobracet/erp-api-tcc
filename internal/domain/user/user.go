package user

import (
	"errors"
	"time"
)

// User representa a entidade de usuário no domínio
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"-"` // Não expor senha no JSON
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

// CreateUserRequest representa os dados necessários para criar um usuário
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required,min=2"`
	Role     string `json:"role" validate:"required,oneof=admin user manager"`
}

// UpdateUserRequest representa os dados para atualizar um usuário
type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=2"`
	Role     *string `json:"role,omitempty" validate:"omitempty,oneof=admin user manager"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// LoginRequest representa os dados para login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse representa a resposta do login
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

// RefreshTokenRequest representa o request para refresh do token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Erros do domínio
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive      = errors.New("user is inactive")
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
)

// ValidateCreate valida os dados para criação de usuário
func (req *CreateUserRequest) ValidateCreate() error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Role == "" {
		return errors.New("role is required")
	}
	
	validRoles := map[string]bool{"admin": true, "user": true, "manager": true}
	if !validRoles[req.Role] {
		return errors.New("invalid role")
	}
	
	return nil
}

// ValidateLogin valida os dados para login
func (req *LoginRequest) ValidateLogin() error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
} 