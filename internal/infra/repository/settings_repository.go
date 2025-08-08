package repository

import (
	"context"
	"errors"

	settingsDomain "erp-api/internal/domain/settings"

	"gorm.io/gorm"
)

type SettingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) settingsDomain.Repository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) Get(ctx context.Context) (*settingsDomain.Settings, error) {
	var settings settingsDomain.Settings
	
	result := r.db.WithContext(ctx).First(&settings)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Se não existir configurações, criar uma padrão
			settings = settingsDomain.Settings{
				TradeName:      "Marmoraria Exemplo",
				LegalName:      "Marmoraria Exemplo LTI",
				CNPJ:           "00.000.000/0000-00",
				Phone:          "(00) 0000-0000",
				Email:          "contato@marmorariaexemplo.com",
				Street:         "Rua Exemplo",
				Number:         "123",
				Complement:     "",
				Neighborhood:   "Centro",
				City:           "São Paulo",
				State:          "SP",
				ZipCode:        "00000-000",
				PrimaryColor:   "#1976d2",
				SecondaryColor: "#9c27b0",
				LogoURL:        "/api/placeholder/200/8C",
			}
			
			result = r.db.WithContext(ctx).Create(&settings)
			if result.Error != nil {
				return nil, result.Error
			}
		} else {
			return nil, result.Error
		}
	}
	
	return &settings, nil
}

func (r *SettingsRepository) Update(ctx context.Context, settings *settingsDomain.Settings) error {
	result := r.db.WithContext(ctx).Save(settings)
	if result.Error != nil {
		return result.Error
	}
	
	return nil
} 