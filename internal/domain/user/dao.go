package user

import "time"

type UserDAO struct {
	ID           string     `db:"id"`
	Email        string     `db:"email"`
	Password     string     `db:"password"`
	Name         string     `db:"name"`
	Role         string     `db:"role"`
	IsActive     bool       `db:"is_active"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	LastLoginAt  *time.Time `db:"last_login_at"`
}

func (dao *UserDAO) ToEntity() *User {
	user := &User{
		ID:        dao.ID,
		Email:     dao.Email,
		Password:  dao.Password,
		Name:      dao.Name,
		Role:      dao.Role,
		IsActive:  dao.IsActive,
		CreatedAt: dao.CreatedAt,
		UpdatedAt: dao.UpdatedAt,
	}
	
	if dao.LastLoginAt != nil {
		user.LastLoginAt = dao.LastLoginAt
	}
	
	return user
}

func (user *User) ToDAO() *UserDAO {
	dao := &UserDAO{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		Name:      user.Name,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	
	if user.LastLoginAt != nil {
		dao.LastLoginAt = user.LastLoginAt
	}
	
	return dao
}

func (user *User) ToDTO() *UserDTO {
	dto := &UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
	
	if user.LastLoginAt != nil {
		dto.LastLoginAt = user.LastLoginAt.Format(time.RFC3339)
	}
	
	return dto
}

func (dto *CreateUserDTO) ToEntity() *CreateUserRequest {
	return &CreateUserRequest{
		Email:    dto.Email,
		Password: dto.Password,
		Name:     dto.Name,
		Role:     dto.Role,
	}
}

func (dto *UpdateUserDTO) ToEntity() *UpdateUserRequest {
	return &UpdateUserRequest{
		Name:     dto.Name,
		Role:     dto.Role,
		IsActive: dto.IsActive,
	}
}

func (dto *LoginDTO) ToEntity() *LoginRequest {
	return &LoginRequest{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func (dto *RefreshTokenDTO) ToEntity() *RefreshTokenRequest {
	return &RefreshTokenRequest{
		RefreshToken: dto.RefreshToken,
	}
} 