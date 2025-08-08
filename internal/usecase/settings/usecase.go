package settings

import (
	"context"
	"time"

	settingsDomain "erp-api/internal/domain/settings"
)

type UseCaseInterface interface {
	Get(ctx context.Context) (*settingsDomain.Settings, error)
	Update(ctx context.Context, req *settingsDomain.UpdateSettingsDTO) (*settingsDomain.Settings, error)
}

type UseCase struct {
	settingsRepo settingsDomain.Repository
}

func NewUseCase(settingsRepo settingsDomain.Repository) UseCaseInterface {
	return &UseCase{
		settingsRepo: settingsRepo,
	}
}

func (u *UseCase) Get(ctx context.Context) (*settingsDomain.Settings, error) {
	return u.settingsRepo.Get(ctx)
}

func (u *UseCase) Update(ctx context.Context, req *settingsDomain.UpdateSettingsDTO) (*settingsDomain.Settings, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Buscar configurações existentes
	settings, err := u.settingsRepo.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Atualizar campos
	if req.TradeName != "" {
		settings.TradeName = req.TradeName
	}
	if req.LegalName != "" {
		settings.LegalName = req.LegalName
	}
	if req.CNPJ != "" {
		settings.CNPJ = req.CNPJ
	}
	if req.Phone != "" {
		settings.Phone = req.Phone
	}
	if req.Email != "" {
		settings.Email = req.Email
	}
	if req.Street != "" {
		settings.Street = req.Street
	}
	if req.Number != "" {
		settings.Number = req.Number
	}
	if req.Complement != "" {
		settings.Complement = req.Complement
	}
	if req.Neighborhood != "" {
		settings.Neighborhood = req.Neighborhood
	}
	if req.City != "" {
		settings.City = req.City
	}
	if req.State != "" {
		settings.State = req.State
	}
	if req.ZipCode != "" {
		settings.ZipCode = req.ZipCode
	}
	if req.PrimaryColor != "" {
		settings.PrimaryColor = req.PrimaryColor
	}
	if req.SecondaryColor != "" {
		settings.SecondaryColor = req.SecondaryColor
	}
	if req.LogoURL != "" {
		settings.LogoURL = req.LogoURL
	}

	settings.UpdatedAt = time.Now()

	err = u.settingsRepo.Update(ctx, settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
} 