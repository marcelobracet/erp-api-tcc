package settings

import "errors"

var (
	ErrSettingsNotFound = errors.New("settings not found")
	ErrInvalidCNPJ      = errors.New("invalid CNPJ")
	ErrInvalidEmail     = errors.New("invalid email")
)

func (req *UpdateSettingsDTO) Validate() error {
	if req.TenantID == "" {
		return errors.New("tenant_id is required")
	}
	if req.Settings == nil {
		return errors.New("settings is required")
	}
	return nil
}

func isValidEmail(email string) bool {
	// Implementação simples de validação de email
	// Em produção, use uma biblioteca como github.com/asaskevich/govalidator
	return len(email) > 0 && len(email) < 255
} 