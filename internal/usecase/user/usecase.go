package user

import (
	"context"
	"time"

	userDomain "erp-api/internal/domain/user"
	"erp-api/internal/utils/dbtypes"
)

// UseCase implementa a lógica de negócio para usuários
type UseCase struct {
	userRepo userDomain.Repository
}

// NewUseCase cria uma nova instância do UseCase
func NewUseCase(userRepo userDomain.Repository) *UseCase {
	return &UseCase{
		userRepo: userRepo,
	}
}

// Register registra um novo usuário
func (u *UseCase) Register(ctx context.Context, req *userDomain.CreateUserRequest) (*userDomain.User, error) {
	// Validar request
	if err := req.ValidateCreate(); err != nil {
		return nil, err
	}

	// Criar usuário como referência a uma identidade externa (Keycloak)
	newUser := &userDomain.User{
		TenantID:    dbtypes.UUID(req.TenantID),
		KeycloakID:  dbtypes.UUID(req.KeycloakID),
		ID:          dbtypes.UUID(req.KeycloakID),
		DisplayName: req.DisplayName,
		Email:       req.Email,
	}

	err := u.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// GetByID busca um usuário pelo ID
func (u *UseCase) GetByID(ctx context.Context, id string) (*userDomain.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

// Update atualiza um usuário
func (u *UseCase) Update(ctx context.Context, id string, req *userDomain.UpdateUserRequest) (*userDomain.User, error) {
	// Buscar usuário existente
	existingUser, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Aplicar atualizações
	if req.DisplayName != nil {
		existingUser.DisplayName = *req.DisplayName
	}
	if req.Email != nil {
		existingUser.Email = req.Email
	}
	
	existingUser.UpdatedAt = time.Now()
	
	// Salvar atualizações
	err = u.userRepo.Update(ctx, existingUser)
	if err != nil {
		return nil, err
	}
	
	return existingUser, nil
}

// Delete remove um usuário
func (u *UseCase) Delete(ctx context.Context, id string) error {
	return u.userRepo.Delete(ctx, id)
}

// List retorna uma lista de usuários
func (u *UseCase) List(ctx context.Context, limit, offset int) ([]*userDomain.User, error) {
	return u.userRepo.List(ctx, limit, offset)
}

// Count retorna o total de usuários
func (u *UseCase) Count(ctx context.Context) (int, error) {
	return u.userRepo.Count(ctx)
}