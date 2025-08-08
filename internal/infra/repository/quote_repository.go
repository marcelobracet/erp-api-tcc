package repository

import (
	"context"
	"errors"

	quoteDomain "erp-api/internal/domain/quote"

	"gorm.io/gorm"
)

type QuoteRepository struct {
	db *gorm.DB
}

func NewQuoteRepository(db *gorm.DB) quoteDomain.Repository {
	return &QuoteRepository{db: db}
}

func (r *QuoteRepository) Create(ctx context.Context, quote *quoteDomain.Quote) error {
	result := r.db.WithContext(ctx).Create(quote)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *QuoteRepository) GetByID(ctx context.Context, id string) (*quoteDomain.Quote, error) {
	var quote quoteDomain.Quote
	
	result := r.db.WithContext(ctx).
		Preload("Client").
		Where("id = ?", id).
		First(&quote)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, quoteDomain.ErrQuoteNotFound
		}
		return nil, result.Error
	}
	
	return &quote, nil
}

func (r *QuoteRepository) Update(ctx context.Context, quote *quoteDomain.Quote) error {
	result := r.db.WithContext(ctx).Save(quote)
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return quoteDomain.ErrQuoteNotFound
	}
	
	return nil
}

func (r *QuoteRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&quoteDomain.Quote{})
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return quoteDomain.ErrQuoteNotFound
	}
	
	return nil
}

func (r *QuoteRepository) List(ctx context.Context, limit, offset int) ([]*quoteDomain.Quote, error) {
	var quotes []*quoteDomain.Quote
	
	result := r.db.WithContext(ctx).
		Preload("Client").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&quotes)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return quotes, nil
}

func (r *QuoteRepository) Count(ctx context.Context) (int, error) {
	var count int64
	
	result := r.db.WithContext(ctx).Model(&quoteDomain.Quote{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	
	return int(count), nil
}

func (r *QuoteRepository) UpdateStatus(ctx context.Context, id string, status quoteDomain.QuoteStatus) error {
	result := r.db.WithContext(ctx).
		Model(&quoteDomain.Quote{}).
		Where("id = ?", id).
		Update("status", status)
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return quoteDomain.ErrQuoteNotFound
	}
	
	return nil
}

type QuoteItemRepository struct {
	db *gorm.DB
}

func NewQuoteItemRepository(db *gorm.DB) quoteDomain.ItemRepository {
	return &QuoteItemRepository{db: db}
}

func (r *QuoteItemRepository) Create(ctx context.Context, item *quoteDomain.QuoteItem) error {
	result := r.db.WithContext(ctx).Create(item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *QuoteItemRepository) GetByQuoteID(ctx context.Context, quoteID string) ([]*quoteDomain.QuoteItem, error) {
	var items []*quoteDomain.QuoteItem
	
	result := r.db.WithContext(ctx).
		Preload("Product").
		Where("quote_id = ?", quoteID).
		Find(&items)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return items, nil
}

func (r *QuoteItemRepository) DeleteByQuoteID(ctx context.Context, quoteID string) error {
	result := r.db.WithContext(ctx).
		Where("quote_id = ?", quoteID).
		Delete(&quoteDomain.QuoteItem{})
	
	if result.Error != nil {
		return result.Error
	}
	
	return nil
} 