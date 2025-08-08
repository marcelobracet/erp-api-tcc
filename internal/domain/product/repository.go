package product

import "context"

type Repository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id string) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*Product, error)
	Count(ctx context.Context) (int, error)
} 