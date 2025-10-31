package client

import "context"

type Repository interface {
	Create(ctx context.Context, client *Client) error
	GetByID(ctx context.Context, tenantID, id string) (*Client, error)
	GetByDocument(ctx context.Context, tenantID, document string) (*Client, error)
	Update(ctx context.Context, client *Client) error
	Delete(ctx context.Context, tenantID, id string) error
	List(ctx context.Context, tenantID string, limit, offset int) ([]*Client, error)
	Count(ctx context.Context, tenantID string) (int, error)
} 