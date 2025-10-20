package repository

import (
	"context"
	"erp-api/internal/domain/tenant"

	"gorm.io/gorm"
)

type tenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) tenant.Repository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) Create(ctx context.Context, tenant *tenant.Tenant) error {
	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *tenantRepository) GetByID(ctx context.Context, id string) (*tenant.Tenant, error) {
	var tenant tenant.Tenant
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepository) Update(ctx context.Context, tenant *tenant.Tenant) error {
	return r.db.WithContext(ctx).Save(tenant).Error
}

func (r *tenantRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&tenant.Tenant{}, "id = ?", id).Error
}

func (r *tenantRepository) List(ctx context.Context, limit, offset int) ([]*tenant.Tenant, error) {
	var tenants []*tenant.Tenant
	err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Offset(offset).Find(&tenants).Error
	return tenants, err
}

func (r *tenantRepository) Count(ctx context.Context) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&tenant.Tenant{}).Count(&count).Error
	return int(count), err
}
