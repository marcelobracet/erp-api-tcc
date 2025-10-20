package settings

type UpdateSettingsDTO struct {
	TenantID string            `json:"tenant_id" binding:"required"`
	Settings map[string]string `json:"settings" binding:"required"`
}

type SettingsDTO struct {
	TenantID string            `json:"tenant_id"`
	Settings map[string]string `json:"settings"`
}