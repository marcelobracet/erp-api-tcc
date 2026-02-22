package repository

import (
	"context"

	settingsDomain "erp-api/internal/domain/settings"
	"erp-api/internal/utils/dbtypes"

	"gorm.io/gorm"
)

type SettingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) settingsDomain.Repository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) Get(ctx context.Context, tenantID string) (map[string]string, error) {
	var settings []settingsDomain.Settings
	
	result := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&settings)
	if result.Error != nil {
		return nil, result.Error
	}

	settingsMap := make(map[string]string)
	for _, setting := range settings {
		settingsMap[setting.Key] = setting.Value
	}

	// Se não existir configurações, retornar configurações padrão
	if len(settingsMap) == 0 {
		settingsMap = map[string]string{
			"company_name":     "Empresa Exemplo",
			"company_email":    "contato@empresa.com",
			"company_phone":    "(11) 3333-4444",
			"company_address":  "Rua Principal, 456",
			"company_city":     "São Paulo",
			"company_state":    "SP",
			"company_zip":      "01234-567",
			"primary_color":    "#2196F3",
			"secondary_color":  "#FFC107",
			"logo_url":         "https://example.com/logo.png",
		}
	}
	
	return settingsMap, nil
}

func (r *SettingsRepository) Update(ctx context.Context, req *settingsDomain.UpdateSettingsDTO) error {
	// Deletar configurações existentes
	err := r.db.WithContext(ctx).Where("tenant_id = ?", req.TenantID).Delete(&settingsDomain.Settings{}).Error
	if err != nil {
		return err
	}

	// Criar novas configurações
	for key, value := range req.Settings {
		setting := &settingsDomain.Settings{
			TenantID: dbtypes.UUID(req.TenantID),
			Key:      key,
			Value:    value,
		}
		
		err := r.db.WithContext(ctx).Create(setting).Error
		if err != nil {
			return err
		}
	}
	
	return nil
}