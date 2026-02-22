package user

import (
	"context"
	userDomain "erp-api/internal/domain/user"
)

// UseCaseInterface define a interface do caso de uso de usuários
type UseCaseInterface interface {
	Register(ctx context.Context, req *userDomain.CreateUserRequest) (*userDomain.User, error)
	GetByID(ctx context.Context, id string) (*userDomain.User, error)
	Update(ctx context.Context, id string, req *userDomain.UpdateUserRequest) (*userDomain.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*userDomain.User, error)
	Count(ctx context.Context) (int, error)
} 