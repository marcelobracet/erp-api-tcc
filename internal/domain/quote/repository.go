package quote

import "context"

type Repository interface {
	Create(ctx context.Context, quote *Quote) error
	GetByID(ctx context.Context, tenantID, id string) (*Quote, error)
	Update(ctx context.Context, quote *Quote) error
	Delete(ctx context.Context, tenantID, id string) error
	List(ctx context.Context, tenantID string, limit, offset int) ([]*Quote, error)
	Count(ctx context.Context, tenantID string) (int, error)
	UpdateStatus(ctx context.Context, tenantID, id string, status QuoteStatus) error
}

type ItemRepository interface {
	Create(ctx context.Context, item *QuoteItem) error
	GetByQuoteID(ctx context.Context, quoteID string) ([]*QuoteItem, error)
	DeleteByQuoteID(ctx context.Context, quoteID string) error
} 