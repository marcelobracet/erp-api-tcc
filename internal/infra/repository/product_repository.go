package repository

import (
	"context"
	"errors"

	productDomain "erp-api/internal/domain/product"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) productDomain.Repository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *productDomain.Product) error {
	result := r.db.WithContext(ctx).Create(product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*productDomain.Product, error) {
	var product productDomain.Product
	
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, productDomain.ErrProductNotFound
		}
		return nil, result.Error
	}
	
	return &product, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *productDomain.Product) error {
	result := r.db.WithContext(ctx).Save(product)
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return productDomain.ErrProductNotFound
	}
	
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&productDomain.Product{})
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return productDomain.ErrProductNotFound
	}
	
	return nil
}

func (r *ProductRepository) List(ctx context.Context, limit, offset int) ([]*productDomain.Product, error) {
	var products []*productDomain.Product
	
	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return products, nil
}

func (r *ProductRepository) Count(ctx context.Context) (int, error) {
	var count int64
	
	result := r.db.WithContext(ctx).Model(&productDomain.Product{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	
	return int(count), nil
} 