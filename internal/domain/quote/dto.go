package quote

type CreateQuoteDTO struct {
	TenantID       string       `json:"tenant_id" binding:"required"`
	ClientID       string       `json:"client_id" binding:"required"`
	UserID         string       `json:"user_id" binding:"required"`
	Discount       float64      `json:"discount,omitempty"`
	Status         QuoteStatus  `json:"status,omitempty"`
	ConversionRate *float64     `json:"conversion_rate,omitempty"`
	Notes          string       `json:"notes,omitempty"`
	Items          []QuoteItemDTO `json:"items" binding:"required"`
}

type UpdateQuoteDTO struct {
	ClientID       string       `json:"client_id,omitempty"`
	UserID         string       `json:"user_id,omitempty"`
	Discount       *float64     `json:"discount,omitempty"`
	Status         QuoteStatus  `json:"status,omitempty"`
	ConversionRate *float64     `json:"conversion_rate,omitempty"`
	Notes          string       `json:"notes,omitempty"`
}

type QuoteItemDTO struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
}

type QuoteDTO struct {
	ID             string       `json:"id"`
	TenantID       string       `json:"tenant_id"`
	ClientID       string       `json:"client_id"`
	UserID         string       `json:"user_id"`
	TotalValue     float64      `json:"total_value"`
	Discount       float64      `json:"discount"`
	Status         QuoteStatus  `json:"status"`
	ConversionRate *float64     `json:"conversion_rate,omitempty"`
	Notes          string       `json:"notes,omitempty"`
	Items          []QuoteItemDTO `json:"items,omitempty"`
	CreatedAt      string       `json:"created_at"`
	UpdatedAt      string       `json:"updated_at"`
}

type QuoteListDTO struct {
	Quotes []*QuoteDTO `json:"quotes"`
	Total  int         `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}

type UpdateQuoteStatusDTO struct {
	Status QuoteStatus `json:"status" binding:"required"`
}