package user

import (
	"context"
	"time"

	userDomain "erp-api/internal/domain/user"
	"erp-api/pkg/auth"
)

// UseCase implementa a lógica de negócio para usuários
type UseCase struct {
	userRepo     userDomain.Repository
	jwtManager   *auth.JWTManager
	passHasher   *auth.PasswordHasher
}

// NewUseCase cria uma nova instância do UseCase
func NewUseCase(userRepo userDomain.Repository, jwtManager *auth.JWTManager, passHasher *auth.PasswordHasher) *UseCase {
	return &UseCase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		passHasher: passHasher,
	}
}

// Register registra um novo usuário
func (u *UseCase) Register(ctx context.Context, req *userDomain.CreateUserRequest) (*userDomain.User, error) {
	// Validar request
	if err := req.ValidateCreate(); err != nil {
		return nil, err
	}
	
	// Verificar se usuário já existe
	existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && err != userDomain.ErrUserNotFound {
		return nil, err
	}
	
	if existingUser != nil {
		return nil, userDomain.ErrUserAlreadyExists
	}
	
	// Hash da senha
	hashedPassword, err := u.passHasher.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	
	// Criar usuário
	newUser := &userDomain.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
		Role:     req.Role,
		IsActive: true,
	}
	
	err = u.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}
	
	return newUser, nil
}

// Login autentica um usuário e retorna tokens
func (u *UseCase) Login(ctx context.Context, req *userDomain.LoginRequest) (*userDomain.LoginResponse, error) {
	// Validar request
	if err := req.ValidateLogin(); err != nil {
		return nil, err
	}
	
	// Buscar usuário por email
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == userDomain.ErrUserNotFound {
			return nil, userDomain.ErrInvalidCredentials
		}
		return nil, err
	}
	
	// Verificar se usuário está ativo
	if !user.IsActive {
		return nil, userDomain.ErrUserInactive
	}
	
	// Verificar senha
	if !u.passHasher.CheckPassword(req.Password, user.Password) {
		return nil, userDomain.ErrInvalidCredentials
	}
	
	// Gerar tokens
	tokenPair, err := u.jwtManager.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	
	// Atualizar último login
	err = u.userRepo.UpdateLastLogin(ctx, user.ID)
	if err != nil {
		// Log do erro mas não falhar o login
		// TODO: implementar logger
	}
	
	return &userDomain.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User:         *user,
	}, nil
}

// RefreshToken gera um novo access token usando refresh token
func (u *UseCase) RefreshToken(ctx context.Context, req *userDomain.RefreshTokenRequest) (*userDomain.LoginResponse, error) {
	// Validar refresh token
	claims, err := u.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, userDomain.ErrInvalidToken
	}
	
	// Buscar usuário
	user, err := u.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	
	// Verificar se usuário está ativo
	if !user.IsActive {
		return nil, userDomain.ErrUserInactive
	}
	
	// Gerar novo access token
	accessToken, err := u.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	
	return &userDomain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken, // Manter o mesmo refresh token
		User:         *user,
	}, nil
}

// GetByID busca um usuário pelo ID
func (u *UseCase) GetByID(ctx context.Context, id string) (*userDomain.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

// GetByEmail busca um usuário pelo email
func (u *UseCase) GetByEmail(ctx context.Context, email string) (*userDomain.User, error) {
	return u.userRepo.GetByEmail(ctx, email)
}

// Update atualiza um usuário
func (u *UseCase) Update(ctx context.Context, id string, req *userDomain.UpdateUserRequest) (*userDomain.User, error) {
	// Buscar usuário existente
	existingUser, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Aplicar atualizações
	if req.Name != nil {
		existingUser.Name = *req.Name
	}
	if req.Role != nil {
		existingUser.Role = *req.Role
	}
	if req.IsActive != nil {
		existingUser.IsActive = *req.IsActive
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

// ValidateToken valida um token JWT e retorna as claims
func (u *UseCase) ValidateToken(tokenString string) (*auth.Claims, error) {
	return u.jwtManager.ValidateToken(tokenString)
} 