package repository

import (
	"context"
	"errors"

	clientDomain "erp-api/internal/domain/client"

	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) clientDomain.Repository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) Create(ctx context.Context, client *clientDomain.Client) error {
	result := r.db.WithContext(ctx).Create(client)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return clientDomain.ErrClientAlreadyExists
		}
		return result.Error
	}
	return nil
}

func (r *ClientRepository) GetByID(ctx context.Context, id string) (*clientDomain.Client, error) {
	var client clientDomain.Client
	
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&client)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, clientDomain.ErrClientNotFound
		}
		return nil, result.Error
	}
	
	return &client, nil
}

func (r *ClientRepository) GetByDocument(ctx context.Context, document string) (*clientDomain.Client, error) {
	var client clientDomain.Client
	
	result := r.db.WithContext(ctx).Where("document = ?", document).First(&client)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, clientDomain.ErrClientNotFound
		}
		return nil, result.Error
	}
	
	return &client, nil
}

func (r *ClientRepository) Update(ctx context.Context, client *clientDomain.Client) error {
	result := r.db.WithContext(ctx).Save(client)
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return clientDomain.ErrClientNotFound
	}
	
	return nil
}

func (r *ClientRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&clientDomain.Client{})
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return clientDomain.ErrClientNotFound
	}
	
	return nil
}

func (r *ClientRepository) List(ctx context.Context, limit, offset int) ([]*clientDomain.Client, error) {
	var clients []*clientDomain.Client
	
	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&clients)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return clients, nil
}

func (r *ClientRepository) Count(ctx context.Context) (int, error) {
	var count int64
	
	result := r.db.WithContext(ctx).Model(&clientDomain.Client{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	
	return int(count), nil
} 