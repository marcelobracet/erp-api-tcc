package quote

import "time"

type CreateQuoteDTO struct {
	ClientID   string       `json:"client_id" binding:"required"`
	Status     QuoteStatus  `json:"status,omitempty"`
	Date       time.Time    `json:"date" binding:"required"`
	ValidUntil time.Time    `json:"valid_until" binding:"required"`
	Notes      string       `json:"notes,omitempty"`
	Items      []QuoteItemDTO `json:"items" binding:"required"`
}

type UpdateQuoteDTO struct {
	ClientID   string       `json:"client_id,omitempty"`
	Status     QuoteStatus  `json:"status,omitempty"`
	Date       *time.Time   `json:"date,omitempty"`
	ValidUntil *time.Time   `json:"valid_until,omitempty"`
	Notes      string       `json:"notes,omitempty"`
	IsActive   *bool        `json:"is_active,omitempty"`
}

type QuoteItemDTO struct {
	ProductID string  `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required"`
}

type QuoteDTO struct {
	ID         string       `json:"id"`
	ClientID   string       `json:"client_id"`
	Client     *Client      `json:"client,omitempty"`
	TotalValue float64      `json:"total_value"`
	Status     QuoteStatus  `json:"status"`
	Date       string       `json:"date"`
	ValidUntil string       `json:"valid_until"`
	Notes      string       `json:"notes,omitempty"`
	IsActive   bool         `json:"is_active"`
	Items      []QuoteItemDTO `json:"items,omitempty"`
	CreatedAt  string       `json:"created_at"`
	UpdatedAt  string       `json:"updated_at"`
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