package product

import "context"

type Repository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, tenantID, id string) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, tenantID, id string) error
	List(ctx context.Context, tenantID string, limit, offset int) ([]*Product, error)
	Count(ctx context.Context, tenantID string) (int, error)
} 