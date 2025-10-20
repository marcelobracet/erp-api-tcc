package user

type UserDTO struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	IsActive     bool   `json:"is_active"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	LastLoginAt  string `json:"last_login_at,omitempty"`
}

type CreateUserDTO struct {
	TenantID string `json:"tenant_id" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required,min=2"`
	Role     string `json:"role" validate:"required,oneof=admin user manager"`
}

type UpdateUserDTO struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=2"`
	Role     *string `json:"role,omitempty" validate:"omitempty,oneof=admin user manager"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDTO struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	User         UserDTO  `json:"user"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UserListDTO struct {
	Users []UserDTO `json:"users"`
	Total int       `json:"total"`
	Limit int       `json:"limit"`
	Offset int      `json:"offset"`
} 