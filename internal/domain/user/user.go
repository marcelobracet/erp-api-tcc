package user

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive      = errors.New("user is inactive")
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
)

func (req *CreateUserRequest) ValidateCreate() error {
	if req.TenantID == "" {
		return errors.New("tenant_id is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Role == "" {
		return errors.New("role is required")
	}
	
	validRoles := map[string]bool{"admin": true, "user": true, "manager": true}
	if !validRoles[req.Role] {
		return errors.New("invalid role")
	}
	
	return nil
}

func (req *LoginRequest) ValidateLogin() error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
} 