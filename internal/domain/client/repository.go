package client

import "context"

type Repository interface {
	Create(ctx context.Context, client *Client) error
	GetByID(ctx context.Context, id string) (*Client, error)
	GetByDocument(ctx context.Context, document string) (*Client, error)
	Update(ctx context.Context, client *Client) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*Client, error)
	Count(ctx context.Context) (int, error)
} 