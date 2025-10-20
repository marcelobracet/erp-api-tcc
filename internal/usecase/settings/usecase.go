package settings

import (
	"context"

	settingsDomain "erp-api/internal/domain/settings"
)

type UseCaseInterface interface {
	Get(ctx context.Context, tenantID string) (*settingsDomain.SettingsDTO, error)
	Update(ctx context.Context, req *settingsDomain.UpdateSettingsDTO) (*settingsDomain.SettingsDTO, error)
}

type UseCase struct {
	settingsRepo settingsDomain.Repository
}

func NewUseCase(settingsRepo settingsDomain.Repository) UseCaseInterface {
	return &UseCase{
		settingsRepo: settingsRepo,
	}
}

func (u *UseCase) Get(ctx context.Context, tenantID string) (*settingsDomain.SettingsDTO, error) {
	settings, err := u.settingsRepo.Get(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	return &settingsDomain.SettingsDTO{
		TenantID:  tenantID,
		Settings: settings,
	}, nil
}

func (u *UseCase) Update(ctx context.Context, req *settingsDomain.UpdateSettingsDTO) (*settingsDomain.SettingsDTO, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	err := u.settingsRepo.Update(ctx, req)
	if err != nil {
		return nil, err
	}

	return &settingsDomain.SettingsDTO{
		TenantID:  req.TenantID,
		Settings: req.Settings,
	}, nil
}