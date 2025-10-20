package user

import (
	"context"
	"testing"
	"time"
)

// MockRepository é uma implementação mock do repositório para testes
type MockRepository struct {
	users map[string]*User
}

// NewMockRepository cria uma nova instância do mock repository
func NewMockRepository() *MockRepository {
	return &MockRepository{
		users: make(map[string]*User),
	}
}

func (m *MockRepository) Create(ctx context.Context, user *User) error {
	if user.ID == "" {
		user.ID = "user-" + time.Now().Format("20060102150405") + "-" + user.Email
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (m *MockRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (m *MockRepository) Update(ctx context.Context, user *User) error {
	if _, exists := m.users[user.ID]; !exists {
		return ErrUserNotFound
	}
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	if _, exists := m.users[id]; !exists {
		return ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *MockRepository) List(ctx context.Context, limit, offset int) ([]*User, error) {
	users := make([]*User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	
	// Simular paginação
	if offset >= len(users) {
		return []*User{}, nil
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
		return ErrUserNotFound
	}
	now := time.Now()
	user.ToDTO().LastLoginAt = now.Format(time.RFC3339)
	user.UpdatedAt = now
	return nil
}

func TestMockRepository_Create(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()
	
	user := &User{
		Email:    "test@example.com",
		Password: "hashed_password",
		Name:     "Test User",
		Role:     "user",
		IsActive: true,
	}
	
	err := repo.Create(ctx, user)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}
	
	if user.ID == "" {
		t.Error("Expected user ID to be set")
	}
	
	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	
	if user.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestMockRepository_GetByID(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()
	
	// Criar um usuário primeiro
	user := &User{
		Email:    "test@example.com",
		Password: "hashed_password",
		Name:     "Test User",
		Role:     "user",
		IsActive: true,
	}
	
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Buscar o usuário criado
	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	
	if found.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, found.Email)
	}
	
	// Testar busca por ID inexistente
	_, err = repo.GetByID(ctx, "non-existent-id")
	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestMockRepository_GetByEmail(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()
	
	// Criar um usuário primeiro
	user := &User{
		Email:    "test@example.com",
		Password: "hashed_password",
		Name:     "Test User",
		Role:     "user",
		IsActive: true,
	}
	
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Buscar o usuário por email
	found, err := repo.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Errorf("GetByEmail() error = %v", err)
	}
	
	if found.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, found.ID)
	}
	
	// Testar busca por email inexistente
	_, err = repo.GetByEmail(ctx, "nonexistent@example.com")
	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestMockRepository_Update(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()
	
	// Criar um usuário primeiro
	user := &User{
		Email:    "test@example.com",
		Password: "hashed_password",
		Name:     "Test User",
		Role:     "user",
		IsActive: true,
	}
	
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Atualizar o usuário
	user.Name = "Updated User"
	err = repo.Update(ctx, user)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}
	
	// Verificar se foi atualizado
	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	
	if found.Name != "Updated User" {
		t.Errorf("Expected name 'Updated User', got '%s'", found.Name)
	}
}

func TestMockRepository_Delete(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()
	
	// Criar um usuário primeiro
	user := &User{
		Email:    "test@example.com",
		Password: "hashed_password",
		Name:     "Test User",
		Role:     "user",
		IsActive: true,
	}
	
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Deletar o usuário
	err = repo.Delete(ctx, user.ID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}
	
	// Verificar se foi deletado
	_, err = repo.GetByID(ctx, user.ID)
	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound after delete, got %v", err)
	}
}

func TestMockRepository_List(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()
	
	// Criar alguns usuários
	users := []*User{
		{Email: "user1@example.com", Password: "pass1", Name: "User 1", Role: "user"},
		{Email: "user2@example.com", Password: "pass2", Name: "User 2", Role: "user"},
		{Email: "user3@example.com", Password: "pass3", Name: "User 3", Role: "user"},
	}
	
	for _, user := range users {
		err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
	}
	
	// Listar usuários
	found, err := repo.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("List() error = %v", err)
	}
	
	if len(found) != 3 {
		t.Errorf("Expected 3 users, got %d", len(found))
	}
	
	// Testar paginação
	found, err = repo.List(ctx, 2, 0)
	if err != nil {
		t.Errorf("List() with limit error = %v", err)
	}
	
	if len(found) != 2 {
		t.Errorf("Expected 2 users with limit, got %d", len(found))
	}
	
	// Testar offset
	found, err = repo.List(ctx, 2, 1)
	if err != nil {
		t.Errorf("List() with offset error = %v", err)
	}
	
	if len(found) != 2 {
		t.Errorf("Expected 2 users with offset, got %d", len(found))
	}
}

func TestMockRepository_Count(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()
	
	// Criar alguns usuários
	users := []*User{
		{Email: "user1@example.com", Password: "pass1", Name: "User 1", Role: "user"},
		{Email: "user2@example.com", Password: "pass2", Name: "User 2", Role: "user"},
	}
	
	for _, user := range users {
		err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
	}
	
	// Contar usuários
	count, err := repo.Count(ctx)
	if err != nil {
		t.Errorf("Count() error = %v", err)
	}
	
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
	
	// Testar count com 0 usuários
	emptyRepo := NewMockRepository()
	count, err = emptyRepo.Count(ctx)
	if err != nil {
		t.Errorf("Count() empty repo error = %v", err)
	}
	
	if count != 0 {
		t.Errorf("Expected count 0 for empty repo, got %d", count)
	}
}

func TestMockRepository_UpdateLastLogin(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()
	
	// Criar um usuário primeiro
	user := &User{
		Email:    "test@example.com",
		Password: "hashed_password",
		Name:     "Test User",
		Role:     "user",
		IsActive: true,
	}
	
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Atualizar último login
	err = repo.UpdateLastLogin(ctx, user.ID)
	if err != nil {
		t.Errorf("UpdateLastLogin() error = %v", err)
	}
	
	// Verificar se foi atualizado
	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	
	if found.ToDTO().LastLoginAt == "" {
		t.Error("Expected LastLoginAt to be set")
	}
} 
