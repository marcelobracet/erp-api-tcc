package product

type CreateProductDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Unit        string  `json:"unit" binding:"required"`
}

type UpdateProductDTO struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Type        string  `json:"type,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Unit        string  `json:"unit,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type ProductDTO struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
	Unit        string  `json:"unit"`
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