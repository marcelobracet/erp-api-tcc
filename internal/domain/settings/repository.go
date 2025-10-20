package settings

import "context"

type Repository interface {
	Get(ctx context.Context, tenantID string) (map[string]string, error)
	Update(ctx context.Context, req *UpdateSettingsDTO) error
} 