package user

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

func (req *CreateUserRequest) ValidateCreate() error {
	if req.TenantID == "" {
		return errors.New("tenant_id is required")
	}
	if req.KeycloakID == "" {
		return errors.New("keycloak_id is required")
	}
	if req.DisplayName == "" {
		return errors.New("display_name is required")	
	}

	return nil
}