package tenant

type CreateTenantDTO struct {
	Name string `json:"name" binding:"required"`
	Plan string `json:"plan,omitempty"`
}

type UpdateTenantDTO struct {
	Name     string `json:"name,omitempty"`
	Plan     string `json:"plan,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}

type TenantDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Plan      string `json:"plan"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TenantListDTO struct {
	Tenants []*TenantDTO `json:"tenants"`
	Total   int          `json:"total"`
	Limit   int          `json:"limit"`
	Offset  int          `json:"offset"`
}
