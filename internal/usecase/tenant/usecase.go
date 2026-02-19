package tenant

import (
	"context"
	"erp-api/internal/domain/tenant"
	"os"
	"strconv"
	"time"
)

const (
	defaultTenantPlanEnvKey   = "TENANT_DEFAULT_PLAN"
	defaultTenantPlanFallback = "free"
	freeTrialDaysEnvKey       = "TENANT_FREE_TRIAL_DAYS"
	freeTrialDaysFallback     = 14
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

func defaultTenantPlan() string {
	plan := os.Getenv(defaultTenantPlanEnvKey)
	if plan == "" {
		return defaultTenantPlanFallback
	}
	return plan
}

func freeTrialDays() int {
	val := os.Getenv(freeTrialDaysEnvKey)
	if val == "" {
		return freeTrialDaysFallback
	}

	days, err := strconv.Atoi(val)
	if err != nil {
		return freeTrialDaysFallback
	}
	if days < 0 {
		return 0
	}
	return days
}

func (u *UseCase) Create(ctx context.Context, dto *tenant.CreateTenantDTO) (*tenant.TenantDTO, error) {
	now := time.Now()
	trialDays := freeTrialDays()
	plan := defaultTenantPlan()

	var trialEndsAt *time.Time
	status := tenant.TenantStatusPending
	if trialDays > 0 {
		t := now.AddDate(0, 0, trialDays)
		trialEndsAt = &t
		status = tenant.TenantStatusActive
	}

	tenantEntity := &tenant.Tenant{
		CompanyName: dto.CompanyName,
		TradeName:   dto.TradeName,
		CNPJ:        dto.CNPJ,
		Email:       dto.Email,
		Phone:       dto.Phone,
		Plan:        plan,
		Status:      status,
		TrialEndsAt: trialEndsAt,
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

	if dto.CompanyName != nil {
		tenantEntity.CompanyName = *dto.CompanyName
	}
	if dto.TradeName != nil {
		tenantEntity.TradeName = *dto.TradeName
	}
	if dto.CNPJ != nil {
		tenantEntity.CNPJ = *dto.CNPJ
	}
	if dto.Email != nil {
		tenantEntity.Email = *dto.Email
	}
	if dto.Phone != nil {
		tenantEntity.Phone = *dto.Phone
	}
	// Plano é único no sistema (configurado via ENV) e não pode ser alterado por request.
	tenantEntity.Plan = defaultTenantPlan()
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
	var trialEndsAt *string
	if entity.TrialEndsAt != nil {
		formatted := entity.TrialEndsAt.Format(time.RFC3339)
		trialEndsAt = &formatted
	}

	return &tenant.TenantDTO{
		ID:          entity.ID,
		CompanyName: entity.CompanyName,
		TradeName:   entity.TradeName,
		CNPJ:        entity.CNPJ,
		Email:       entity.Email,
		Phone:       entity.Phone,
		Plan:        entity.Plan,
		Status:      string(entity.Status),
		IsActive:    entity.IsActive,
		TrialEndsAt: trialEndsAt,
		CreatedAt:   entity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   entity.UpdatedAt.Format(time.RFC3339),
	}
}
