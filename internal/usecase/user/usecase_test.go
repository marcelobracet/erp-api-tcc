package user

import (
	"context"
	"testing"
	"time"

	userDomain "erp-api/internal/domain/user"
	"erp-api/pkg/auth"
)

// MockRepository é um mock do repositório para testes
type MockRepository struct {
	users map[string]*userDomain.User
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		users: make(map[string]*userDomain.User),
	}
}

func (m *MockRepository) Create(ctx context.Context, user *userDomain.User) error {
	if user.ID == "" {
		user.ID = "user-" + time.Now().Format("20060102150405") + "-" + user.Email
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*userDomain.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, userDomain.ErrUserNotFound
	}
	return user, nil
}

func (m *MockRepository) GetByEmail(ctx context.Context, email string) (*userDomain.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, userDomain.ErrUserNotFound
}

func (m *MockRepository) Update(ctx context.Context, user *userDomain.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return userDomain.ErrUserNotFound
	}
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	if _, exists := m.users[id]; !exists {
		return userDomain.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *MockRepository) List(ctx context.Context, limit, offset int) ([]*userDomain.User, error) {
	users := make([]*userDomain.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}

	if offset >= len(users) {
		return []*userDomain.User{}, nil
	}

	end := offset + limit
	if end > len(users) {
		end = len(users)
	}

	return users[offset:end], nil
}

func (m *MockRepository) Count(ctx context.Context) (int, error) {
	return len(m.users), nil
}

func (m *MockRepository) UpdateLastLogin(ctx context.Context, id string) error {
	user, exists := m.users[id]
	if !exists {
		return userDomain.ErrUserNotFound
	}
	now := time.Now()
	user.LastLoginAt = &now
	user.UpdatedAt = now
	return nil
}

func TestUseCase_Register(t *testing.T) {
	repo := NewMockRepository()
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	passHasher := auth.DefaultPasswordHasher()

	useCase := NewUseCase(repo, jwtManager, passHasher)
	ctx := context.Background()

	req := &userDomain.CreateUserRequest{
		TenantID: "tenant-123",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "user",
	}

	user, err := useCase.Register(ctx, req)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if user.Email != req.Email {
		t.Errorf("Expected email %s, got %s", req.Email, user.Email)
	}

	if user.Name != req.Name {
		t.Errorf("Expected name %s, got %s", req.Name, user.Name)
	}

	if user.Role != req.Role {
		t.Errorf("Expected role %s, got %s", req.Role, user.Role)
	}

	if !user.IsActive {
		t.Error("Expected user to be active")
	}

	// Verificar se senha foi hasheada
	if user.Password == req.Password {
		t.Error("Expected password to be hashed")
	}

	// Verificar se usuário foi salvo no repositório
	found, err := repo.GetByEmail(ctx, req.Email)
	if err != nil {
		t.Fatalf("GetByEmail() error = %v", err)
	}

	if found.ID != user.ID {
		t.Errorf("Expected user ID %s, got %s", user.ID, found.ID)
	}
}

func TestUseCase_Register_DuplicateEmail(t *testing.T) {
	repo := NewMockRepository()
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	passHasher := auth.DefaultPasswordHasher()

	useCase := NewUseCase(repo, jwtManager, passHasher)
	ctx := context.Background()

	// Criar primeiro usuário
	req1 := &userDomain.CreateUserRequest{
		TenantID: "tenant-123",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "user",
	}

	_, err := useCase.Register(ctx, req1)
	if err != nil {
		t.Fatalf("First Register() error = %v", err)
	}

	// Tentar criar segundo usuário com mesmo email
	req2 := &userDomain.CreateUserRequest{
		TenantID: "tenant-123",
		Email:    "test@example.com",
		Password: "password456",
		Name:     "Another User",
		Role:     "user",
	}

	_, err = useCase.Register(ctx, req2)
	if err != userDomain.ErrUserAlreadyExists {
		t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestUseCase_Login(t *testing.T) {
	repo := NewMockRepository()
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	passHasher := auth.DefaultPasswordHasher()

	useCase := NewUseCase(repo, jwtManager, passHasher)
	ctx := context.Background()

	// Criar usuário primeiro
	createReq := &userDomain.CreateUserRequest{
		TenantID: "tenant-123",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "user",
	}

	_, err := useCase.Register(ctx, createReq)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Fazer login
	loginReq := &userDomain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	loginResp, err := useCase.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	if loginResp.AccessToken == "" {
		t.Error("Expected access token to be generated")
	}

	if loginResp.RefreshToken == "" {
		t.Error("Expected refresh token to be generated")
	}

	if loginResp.User.Email != loginReq.Email {
		t.Errorf("Expected user email %s, got %s", loginReq.Email, loginResp.User.Email)
	}
}

func TestUseCase_Login_InvalidCredentials(t *testing.T) {
	repo := NewMockRepository()
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	passHasher := auth.DefaultPasswordHasher()

	useCase := NewUseCase(repo, jwtManager, passHasher)
	ctx := context.Background()

	// Criar usuário primeiro
	createReq := &userDomain.CreateUserRequest{
		TenantID: "tenant-123",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "user",
	}

	_, err := useCase.Register(ctx, createReq)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Tentar login com senha incorreta
	loginReq := &userDomain.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, err = useCase.Login(ctx, loginReq)
	if err != userDomain.ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestUseCase_Login_UserNotFound(t *testing.T) {
	repo := NewMockRepository()
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	passHasher := auth.DefaultPasswordHasher()

	useCase := NewUseCase(repo, jwtManager, passHasher)
	ctx := context.Background()

	// Tentar login com usuário inexistente
	loginReq := &userDomain.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	_, err := useCase.Login(ctx, loginReq)
	if err != userDomain.ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestUseCase_RefreshToken(t *testing.T) {
	repo := NewMockRepository()
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	passHasher := auth.DefaultPasswordHasher()

	useCase := NewUseCase(repo, jwtManager, passHasher)
	ctx := context.Background()

	// Criar usuário e fazer login
	createReq := &userDomain.CreateUserRequest{
		TenantID: "tenant-123",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "user",
	}

	_, err := useCase.Register(ctx, createReq)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	loginReq := &userDomain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	loginResp, err := useCase.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	// Refresh token
	refreshReq := &userDomain.RefreshTokenRequest{
		RefreshToken: loginResp.RefreshToken,
	}

	refreshResp, err := useCase.RefreshToken(ctx, refreshReq)
	if err != nil {
		t.Fatalf("RefreshToken() error = %v", err)
	}

	if refreshResp.AccessToken == "" {
		t.Error("Expected new access token to be generated")
	}

	if refreshResp.RefreshToken != loginResp.RefreshToken {
		t.Error("Expected refresh token to remain the same")
	}

	if refreshResp.User.Email != loginResp.User.Email {
		t.Errorf("Expected user email %s, got %s", loginResp.User.Email, refreshResp.User.Email)
	}
}

func TestUseCase_GetByID(t *testing.T) {
	repo := NewMockRepository()
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	passHasher := auth.DefaultPasswordHasher()

	useCase := NewUseCase(repo, jwtManager, passHasher)
	ctx := context.Background()

	// Criar usuário
	createReq := &userDomain.CreateUserRequest{
		TenantID: "tenant-123",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "user",
	}

	createdUser, err := useCase.Register(ctx, createReq)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Buscar usuário por ID
	foundUser, err := useCase.GetByID(ctx, createdUser.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if foundUser.ID != createdUser.ID {
		t.Errorf("Expected user ID %s, got %s", createdUser.ID, foundUser.ID)
	}

	if foundUser.Email != createdUser.Email {
		t.Errorf("Expected user email %s, got %s", createdUser.Email, foundUser.Email)
	}
}
