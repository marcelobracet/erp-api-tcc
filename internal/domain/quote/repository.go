package quote

import "context"

type Repository interface {
	Create(ctx context.Context, quote *Quote) error
	GetByID(ctx context.Context, id string) (*Quote, error)
	Update(ctx context.Context, quote *Quote) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*Quote, error)
	Count(ctx context.Context) (int, error)
	UpdateStatus(ctx context.Context, id string, status QuoteStatus) error
}

type ItemRepository interface {
	Create(ctx context.Context, item *QuoteItem) error
	GetByQuoteID(ctx context.Context, quoteID string) ([]*QuoteItem, error)
	DeleteByQuoteID(ctx context.Context, quoteID string) error
} 