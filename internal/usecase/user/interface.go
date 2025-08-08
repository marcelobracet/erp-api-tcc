package user

import (
	"context"
	userDomain "erp-api/internal/domain/user"
	"erp-api/pkg/auth"
)

// UseCaseInterface define a interface do caso de uso de usuários
type UseCaseInterface interface {
	Register(ctx context.Context, req *userDomain.CreateUserRequest) (*userDomain.User, error)
	Login(ctx context.Context, req *userDomain.LoginRequest) (*userDomain.LoginResponse, error)
	RefreshToken(ctx context.Context, req *userDomain.RefreshTokenRequest) (*userDomain.LoginResponse, error)
	GetByID(ctx context.Context, id string) (*userDomain.User, error)
	GetByEmail(ctx context.Context, email string) (*userDomain.User, error)
	Update(ctx context.Context, id string, req *userDomain.UpdateUserRequest) (*userDomain.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*userDomain.User, error)
	Count(ctx context.Context) (int, error)
	ValidateToken(tokenString string) (*auth.Claims, error)
} 