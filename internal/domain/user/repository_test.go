package user

import (
	"context"
	"testing"
	"time"

	"erp-api/internal/utils/dbtypes"
)

type MockRepository struct {
	users map[string]*User
}

func strPtr(s string) *string { return &s }

func NewMockRepository() *MockRepository {
	return &MockRepository{
		users: make(map[string]*User),
	}
}

func (m *MockRepository) Create(ctx context.Context, user *User) error {
	if user.ID == "" {
		user.ID = dbtypes.NewUUID()
	}
	if user.KeycloakID == "" {
		user.KeycloakID = user.ID
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[string(user.ID)] = user
	return nil
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (m *MockRepository) Update(ctx context.Context, user *User) error {
	if _, exists := m.users[string(user.ID)]; !exists {
		return ErrUserNotFound
	}
	user.UpdatedAt = time.Now()
	m.users[string(user.ID)] = user
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

func TestMockRepository_Create(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()

	user := &User{
		TenantID:    "tenant-123",
		KeycloakID:  "kc-123",
		DisplayName: "Test User",
		Email:       strPtr("test@example.com"),
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	if user.ID == "" {
		t.Error("Expected user ID to be set")
	}
	if user.KeycloakID == "" {
		t.Error("Expected KeycloakID to be set")
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

	user := &User{
		TenantID:    "tenant-123",
		KeycloakID:  "kc-123",
		DisplayName: "Test User",
		Email:       strPtr("test@example.com"),
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	found, err := repo.GetByID(ctx, string(user.ID))
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if string(found.ID) != string(user.ID) {
		t.Errorf("Expected ID %s, got %s", string(user.ID), string(found.ID))
	}

	_, err = repo.GetByID(ctx, "non-existent-id")
	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestMockRepository_Update(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()

	user := &User{
		TenantID:    "tenant-123",
		KeycloakID:  "kc-123",
		DisplayName: "Test User",
		Email:       strPtr("test@example.com"),
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user.DisplayName = "Updated User"
	err = repo.Update(ctx, user)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	found, err := repo.GetByID(ctx, string(user.ID))
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}

	if found.DisplayName != "Updated User" {
		t.Errorf("Expected display_name 'Updated User', got '%s'", found.DisplayName)
	}
}

func TestMockRepository_Delete(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()

	user := &User{
		TenantID:    "tenant-123",
		KeycloakID:  "kc-123",
		DisplayName: "Test User",
		Email:       strPtr("test@example.com"),
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	err = repo.Delete(ctx, string(user.ID))
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	_, err = repo.GetByID(ctx, string(user.ID))
	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound after delete, got %v", err)
	}
}

func TestMockRepository_List(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()

	users := []*User{
		{TenantID: "tenant-123", KeycloakID: "kc-1", DisplayName: "User 1", Email: strPtr("user1@example.com")},
		{TenantID: "tenant-123", KeycloakID: "kc-2", DisplayName: "User 2", Email: strPtr("user2@example.com")},
		{TenantID: "tenant-123", KeycloakID: "kc-3", DisplayName: "User 3", Email: strPtr("user3@example.com")},
	}

	for _, user := range users {
		err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
	}

	found, err := repo.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("List() error = %v", err)
	}

	if len(found) != 3 {
		t.Errorf("Expected 3 users, got %d", len(found))
	}

	found, err = repo.List(ctx, 2, 0)
	if err != nil {
		t.Errorf("List() with limit error = %v", err)
	}

	if len(found) != 2 {
		t.Errorf("Expected 2 users with limit, got %d", len(found))
	}

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

	users := []*User{
		{TenantID: "tenant-123", KeycloakID: "kc-1", DisplayName: "User 1", Email: strPtr("user1@example.com")},
		{TenantID: "tenant-123", KeycloakID: "kc-2", DisplayName: "User 2", Email: strPtr("user2@example.com")},
	}

	for _, user := range users {
		err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
	}

	count, err := repo.Count(ctx)
	if err != nil {
		t.Errorf("Count() error = %v", err)
	}

	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}

	emptyRepo := NewMockRepository()
	count, err = emptyRepo.Count(ctx)
	if err != nil {
		t.Errorf("Count() empty repo error = %v", err)
	}

	if count != 0 {
		t.Errorf("Expected count 0 for empty repo, got %d", count)
	}
}
