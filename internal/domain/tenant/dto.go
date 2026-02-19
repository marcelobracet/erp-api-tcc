package tenant

type CreateTenantDTO struct {
	CompanyName string `json:"company_name" binding:"required,min=2"`
	TradeName   string `json:"trade_name,omitempty"`
	CNPJ        string `json:"cnpj,omitempty"`
	Email       string `json:"email,omitempty" binding:"omitempty,email"`
	Phone       string `json:"phone,omitempty"`

	Plan string `json:"plan,omitempty"` // default: free
}

type UpdateTenantDTO struct {
	CompanyName *string `json:"company_name,omitempty"`
	TradeName   *string `json:"trade_name,omitempty"`
	CNPJ        *string `json:"cnpj,omitempty"`
	Email       *string `json:"email,omitempty" binding:"omitempty,email"`
	Phone       *string `json:"phone,omitempty"`

	Plan     *string `json:"plan,omitempty"`
	Status   *string `json:"status,omitempty"` // active | suspended | canceled
	IsActive *bool   `json:"is_active,omitempty"`
}
type TenantDTO struct {
	ID string `json:"id"`

	CompanyName string `json:"company_name"`
	TradeName   string `json:"trade_name,omitempty"`
	CNPJ        string `json:"cnpj,omitempty"`

	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`

	Plan        string  `json:"plan"`
	Status      string  `json:"status"`
	IsActive    bool    `json:"is_active"`
	TrialEndsAt *string `json:"trial_ends_at,omitempty"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TenantListDTO struct {
	Tenants []*TenantDTO `json:"tenants"`
	Total   int          `json:"total"`
	Limit   int          `json:"limit"`
	Offset  int          `json:"offset"`
}
