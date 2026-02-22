package user

type UserDTO struct {
	ID          string  `json:"id"`
	KeycloakID  string  `json:"keycloak_id"`
	TenantID    string  `json:"tenant_id"`
	DisplayName string  `json:"display_name"`
	Email       *string `json:"email,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type CreateUserDTO struct {
	TenantID    string  `json:"tenant_id" validate:"required"`
	KeycloakID  string  `json:"keycloak_id" validate:"required"`
	DisplayName string  `json:"display_name" validate:"required,min=2"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
}

type UpdateUserDTO struct {
	DisplayName *string `json:"display_name,omitempty" validate:"omitempty,min=2"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
}

type UserListDTO struct {
	Users  []UserDTO `json:"users"`
	Total  int       `json:"total"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}
