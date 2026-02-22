package user

import (
	"time"

	"erp-api/internal/utils/dbtypes"
)

type UserDAO struct {
	ID         string  `db:"id"`
	KeycloakID string  `db:"keycloak_id"`
	TenantID   string  `db:"tenant_id"`
	Email      *string `db:"email"`
	Name       string  `db:"display_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (dao *UserDAO) ToEntity() *User {
	user := &User{
		ID:          dbtypes.UUID(dao.ID),
		KeycloakID:  dbtypes.UUID(dao.KeycloakID),
		TenantID:    dbtypes.UUID(dao.TenantID),
		Email:       dao.Email,
		DisplayName: dao.Name,
		CreatedAt: dao.CreatedAt,
		UpdatedAt: dao.UpdatedAt,
	}

	return user
}

func (user *User) ToDAO() *UserDAO {
	dao := &UserDAO{
		ID:         string(user.ID),
		KeycloakID: string(user.KeycloakID),
		TenantID:   string(user.TenantID),
		Email:      user.Email,
		Name:       user.DisplayName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return dao
}

func (user *User) ToDTO() *UserDTO {
	dto := &UserDTO{
		ID:          string(user.ID),
		KeycloakID:  string(user.KeycloakID),
		TenantID:    string(user.TenantID),
		DisplayName: user.DisplayName,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}

	return dto
}

func (dto *CreateUserDTO) ToEntity() *CreateUserRequest {
	return &CreateUserRequest{
		TenantID:    dto.TenantID,
		KeycloakID:  dto.KeycloakID,
		DisplayName: dto.DisplayName,
		Email:       dto.Email,
	}
}

func (dto *UpdateUserDTO) ToEntity() *UpdateUserRequest {
	return &UpdateUserRequest{
		DisplayName: dto.DisplayName,
		Email:       dto.Email,
	}
}
