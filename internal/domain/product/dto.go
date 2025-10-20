package product

type CreateProductDTO struct {
	TenantID    string  `json:"tenant_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int     `json:"stock,omitempty"`
	SKU         string  `json:"sku,omitempty"`
	Category    string  `json:"category,omitempty"`
	ImageURL    string  `json:"image_url,omitempty"`
}

type UpdateProductDTO struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Stock       *int    `json:"stock,omitempty"`
	SKU         string  `json:"sku,omitempty"`
	Category    string  `json:"category,omitempty"`
	ImageURL    string  `json:"image_url,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type ProductDTO struct {
	ID          string  `json:"id"`
	TenantID    string  `json:"tenant_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	SKU         string  `json:"sku"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ProductListDTO struct {
	Products []*ProductDTO `json:"products"`
	Total    int           `json:"total"`
	Limit    int           `json:"limit"`
	Offset   int           `json:"offset"`
} 