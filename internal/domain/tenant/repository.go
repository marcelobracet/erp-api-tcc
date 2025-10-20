package tenant

import "context"

type Repository interface {
	Create(ctx context.Context, tenant *Tenant) error
	GetByID(ctx context.Context, id string) (*Tenant, error)
	Update(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*Tenant, error)
	Count(ctx context.Context) (int, error)
}
