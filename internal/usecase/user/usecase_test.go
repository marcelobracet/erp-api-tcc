package user

import (
	"context"
	"testing"
	"time"

	userDomain "erp-api/internal/domain/user"
)

type MockRepository struct {
	users map[string]*userDomain.User
}

func NewMockRepository() *MockRepository {
	return &MockRepository{users: make(map[string]*userDomain.User)}
}

func (m *MockRepository) Create(ctx context.Context, user *userDomain.User) error {
	if _, exists := m.users[string(user.ID)]; exists {
		return userDomain.ErrUserAlreadyExists
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[string(user.ID)] = user
	return nil
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*userDomain.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, userDomain.ErrUserNotFound
	}
	return u, nil
}

func (m *MockRepository) Update(ctx context.Context, user *userDomain.User) error {
	if _, ok := m.users[string(user.ID)]; !ok {
		return userDomain.ErrUserNotFound
	}
	user.UpdatedAt = time.Now()
	m.users[string(user.ID)] = user
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.users[id]; !ok {
		return userDomain.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *MockRepository) List(ctx context.Context, limit, offset int) ([]*userDomain.User, error) {
	all := make([]*userDomain.User, 0, len(m.users))
	for _, u := range m.users {
		all = append(all, u)
	}
	if offset >= len(all) {
		return []*userDomain.User{}, nil
	}
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}
	return all[offset:end], nil
}

func (m *MockRepository) Count(ctx context.Context) (int, error) {
	return len(m.users), nil
}

func strPtr(s string) *string { return &s }

func TestUseCase_Register(t *testing.T) {
	repo := NewMockRepository()
	useCase := NewUseCase(repo)
	ctx := context.Background()

	email := "test@example.com"
	req := &userDomain.CreateUserRequest{
		TenantID:    "tenant-123",
		KeycloakID:  "kc-123",
		DisplayName: "Test User",
		Email:       &email,
	}

	u, err := useCase.Register(ctx, req)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if string(u.ID) != req.KeycloakID {
		t.Fatalf("expected ID=%q, got %q", req.KeycloakID, string(u.ID))
	}
	if string(u.KeycloakID) != req.KeycloakID {
		t.Fatalf("expected KeycloakID=%q, got %q", req.KeycloakID, string(u.KeycloakID))
	}
	if u.DisplayName != req.DisplayName {
		t.Fatalf("expected display_name=%q, got %q", req.DisplayName, u.DisplayName)
	}
	if u.Email == nil || *u.Email != email {
		t.Fatalf("expected email=%q, got %v", email, u.Email)
	}
}

func TestUseCase_Register_Duplicate(t *testing.T) {
	repo := NewMockRepository()
	useCase := NewUseCase(repo)
	ctx := context.Background()

	req := &userDomain.CreateUserRequest{TenantID: "tenant-123", KeycloakID: "kc-123", DisplayName: "User"}
	_, err := useCase.Register(ctx, req)
	if err != nil {
		t.Fatalf("first Register() error = %v", err)
	}

	_, err = useCase.Register(ctx, req)
	if err != userDomain.ErrUserAlreadyExists {
		t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestUseCase_Update(t *testing.T) {
	repo := NewMockRepository()
	useCase := NewUseCase(repo)
	ctx := context.Background()

	_, err := useCase.Register(ctx, &userDomain.CreateUserRequest{TenantID: "tenant-123", KeycloakID: "kc-123", DisplayName: "Old"})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	newName := "New"
	newEmail := "new@example.com"
	u, err := useCase.Update(ctx, "kc-123", &userDomain.UpdateUserRequest{DisplayName: &newName, Email: &newEmail})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if u.DisplayName != newName {
		t.Fatalf("expected display_name=%q, got %q", newName, u.DisplayName)
	}
	if u.Email == nil || *u.Email != newEmail {
		t.Fatalf("expected email=%q, got %v", newEmail, u.Email)
	}
}

func TestUseCase_Delete(t *testing.T) {
	repo := NewMockRepository()
	useCase := NewUseCase(repo)
	ctx := context.Background()

	_, err := useCase.Register(ctx, &userDomain.CreateUserRequest{TenantID: "tenant-123", KeycloakID: "kc-123", DisplayName: "User"})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if err := useCase.Delete(ctx, "kc-123"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = useCase.GetByID(ctx, "kc-123")
	if err != userDomain.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUseCase_ListAndCount(t *testing.T) {
	repo := NewMockRepository()
	useCase := NewUseCase(repo)
	ctx := context.Background()

	_, _ = useCase.Register(ctx, &userDomain.CreateUserRequest{TenantID: "tenant-123", KeycloakID: "kc-1", DisplayName: "U1"})
	_, _ = useCase.Register(ctx, &userDomain.CreateUserRequest{TenantID: "tenant-123", KeycloakID: "kc-2", DisplayName: "U2"})

	count, err := useCase.Count(ctx)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}
	if count != 2 {
		t.Fatalf("expected count=2, got %d", count)
	}

	users, err := useCase.List(ctx, 10, 0)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}
