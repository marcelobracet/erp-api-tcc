package client

type CreateClientDTO struct {
	TenantID     string `json:"tenant_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Document     string `json:"document" binding:"required"`
	DocumentType string `json:"document_type" binding:"required"`
	Address      string `json:"address,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	ZipCode      string `json:"zip_code,omitempty"`
}

type UpdateClientDTO struct {
	Name         string `json:"name,omitempty"`
	Email        string `json:"email,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Document     string `json:"document,omitempty"`
	DocumentType string `json:"document_type,omitempty"`
	Address      string `json:"address,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	ZipCode      string `json:"zip_code,omitempty"`
	IsActive     *bool  `json:"is_active,omitempty"`
}

type ClientDTO struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Document     string `json:"document"`
	DocumentType string `json:"document_type"`
	Address      string `json:"address,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	ZipCode      string `json:"zip_code,omitempty"`
	IsActive     bool   `json:"is_active"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ClientListDTO struct {
	Clients []*ClientDTO `json:"clients"`
	Total   int          `json:"total"`
	Limit   int          `json:"limit"`
	Offset  int          `json:"offset"`
} 