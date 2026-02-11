package tenant

import (
	"context"
	"erp-api/internal/domain/tenant"
	"time"
)

type UseCase struct {
	repo tenant.Repository
}

type UseCaseInterface interface {
	Create(ctx context.Context, dto *tenant.CreateTenantDTO) (*tenant.TenantDTO, error)
	GetByID(ctx context.Context, id string) (*tenant.TenantDTO, error)
	Update(ctx context.Context, id string, dto *tenant.UpdateTenantDTO) (*tenant.TenantDTO, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) (*tenant.TenantListDTO, error)
	Count(ctx context.Context) (int, error)
}

func NewUseCase(repo tenant.Repository) UseCaseInterface {
	return &UseCase{repo: repo}
}

func (u *UseCase) Create(ctx context.Context, dto *tenant.CreateTenantDTO) (*tenant.TenantDTO, error) {
	tenantEntity := &tenant.Tenant{
		CompanyName: dto.CompanyName,
		Plan:        dto.Plan,
		IsActive:    true,
	}

	if err := u.repo.Create(ctx, tenantEntity); err != nil {
		return nil, err
	}

	return u.entityToDTO(tenantEntity), nil
}

func (u *UseCase) GetByID(ctx context.Context, id string) (*tenant.TenantDTO, error) {
	tenantEntity, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return u.entityToDTO(tenantEntity), nil
}

func (u *UseCase) Update(ctx context.Context, id string, dto *tenant.UpdateTenantDTO) (*tenant.TenantDTO, error) {
	tenantEntity, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if *dto.CompanyName != "" {
		tenantEntity.CompanyName = *dto.CompanyName
	}
	if dto.Plan != nil && *dto.Plan != "" {
		tenantEntity.Plan = *dto.Plan
	}
	if dto.IsActive != nil {
		tenantEntity.IsActive = *dto.IsActive
	}

	if err := u.repo.Update(ctx, tenantEntity); err != nil {
		return nil, err
	}

	return u.entityToDTO(tenantEntity), nil
}

func (u *UseCase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *UseCase) List(ctx context.Context, limit, offset int) (*tenant.TenantListDTO, error) {
	tenants, err := u.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := u.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	tenantDTOs := make([]*tenant.TenantDTO, len(tenants))
	for i, t := range tenants {
		tenantDTOs[i] = u.entityToDTO(t)
	}

	return &tenant.TenantListDTO{
		Tenants: tenantDTOs,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	}, nil
}

func (u *UseCase) Count(ctx context.Context) (int, error) {
	return u.repo.Count(ctx)
}

func (u *UseCase) entityToDTO(entity *tenant.Tenant) *tenant.TenantDTO {
	return &tenant.TenantDTO{
		ID:          entity.ID,
		CompanyName: entity.CompanyName,
		Plan:        entity.Plan,
		IsActive:    entity.IsActive,
		CreatedAt:   entity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   entity.UpdatedAt.Format(time.RFC3339),
	}
}
