//go:build legacy_user_tests
// +build legacy_user_tests

package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	userDomain "erp-api/internal/domain/user"
	"erp-api/internal/utils/dbtypes"
	userUseCase "erp-api/internal/usecase/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MockUseCase struct {
	users map[string]*userDomain.User
}

var _ userUseCase.UseCaseInterface = (*MockUseCase)(nil)

func NewMockUseCase() *MockUseCase {
	return &MockUseCase{users: make(map[string]*userDomain.User)}
}

func (m *MockUseCase) Register(ctx context.Context, req *userDomain.CreateUserRequest) (*userDomain.User, error) {
	if err := req.ValidateCreate(); err != nil {
		return nil, err
	}

	// enforce unique keycloak_id

	router := gin.New()
	router.POST("/login", handler.Login)

	// Teste de login válido
	loginReq := userDomain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	jsonBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verificar se retornou tokens
	var response userDomain.LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.AccessToken == "" {
		t.Error("Expected access token to be present")
	}

	if response.RefreshToken == "" {
		t.Error("Expected refresh token to be present")
	}
}

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := NewMockUseCase()
	handler := NewHandler(mockUseCase)

	router := gin.New()
	router.POST("/login", handler.Login)

	// Teste de login com credenciais inválidas
	loginReq := userDomain.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "wrongpassword",
	}

	jsonBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := NewMockUseCase()
	handler := NewHandler(mockUseCase)

	// Criar usuário primeiro
	createReq := userDomain.CreateUserRequest{
		TenantID: "tenant-123",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "user",
	}

	user, err := mockUseCase.Register(context.Background(), &createReq)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	router := gin.New()
	router.GET("/users/:id", handler.GetByID)

	// Teste de busca por ID válido
	req := httptest.NewRequest("GET", "/users/"+user.ID, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Teste de busca por ID inexistente
	req = httptest.NewRequest("GET", "/users/nonexistent-id", nil)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := NewMockUseCase()
	handler := NewHandler(mockUseCase)

	// Criar usuário primeiro
	createReq := userDomain.CreateUserRequest{
		TenantID: "tenant-123",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "user",
	}

	user, err := mockUseCase.Register(context.Background(), &createReq)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	router := gin.New()
	router.PUT("/users/:id", handler.Update)

	// Teste de atualização válida
	newName := "Updated User"
	updateReq := userDomain.UpdateUserRequest{
		Name: &newName,
	}

	jsonBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/users/"+user.ID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verificar se o nome foi atualizado
	var updatedUser userDomain.User
	if err := json.Unmarshal(w.Body.Bytes(), &updatedUser); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if updatedUser.Name != newName {
		t.Errorf("Expected name %s, got %s", newName, updatedUser.Name)
	}
}

func TestHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := NewMockUseCase()
	handler := NewHandler(mockUseCase)

	// Criar usuário primeiro
	createReq := userDomain.CreateUserRequest{
		TenantID: "tenant-123",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "user",
	}

	user, err := mockUseCase.Register(context.Background(), &createReq)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	router := gin.New()
	router.DELETE("/users/:id", handler.Delete)

	// Teste de deleção válida
	req := httptest.NewRequest("DELETE", "/users/"+user.ID, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Verificar se o usuário foi deletado
	_, err = mockUseCase.GetByID(context.Background(), user.ID)
	if err != userDomain.ErrUserNotFound {
		t.Error("Expected user to be deleted")
	}
}

func TestHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := NewMockUseCase()
	handler := NewHandler(mockUseCase)

	// Criar alguns usuários
	users := []userDomain.CreateUserRequest{
		{TenantID: "tenant-123", Email: "user1@example.com", Password: "password1", Name: "User 1", Role: "user"},
		{TenantID: "tenant-123", Email: "user2@example.com", Password: "password2", Name: "User 2", Role: "user"},
		{TenantID: "tenant-123", Email: "user3@example.com", Password: "password3", Name: "User 3", Role: "user"},
	}

	for _, createReq := range users {
		_, err := mockUseCase.Register(context.Background(), &createReq)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
	}

	router := gin.New()
	router.GET("/users", handler.List)

	// Teste de listagem
	req := httptest.NewRequest("GET", "/users", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verificar se retornou a lista de usuários
	var response []*userDomain.User
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 3 {
		t.Errorf("Expected 3 users, got %d", len(response))
	}

	// Verificar se os usuários foram criados corretamente
	count, err := mockUseCase.Count(context.Background())
	if err != nil {
		t.Fatalf("Failed to count users: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 users in mock, got %d", count)
	}
}

func TestHandler_Count(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := NewMockUseCase()
	handler := NewHandler(mockUseCase)

	// Criar alguns usuários
	users := []userDomain.CreateUserRequest{
		{TenantID: "tenant-123", Email: "user1@example.com", Password: "password1", Name: "User 1", Role: "user"},
		{TenantID: "tenant-123", Email: "user2@example.com", Password: "password2", Name: "User 2", Role: "user"},
	}

	for _, createReq := range users {
		_, err := mockUseCase.Register(context.Background(), &createReq)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
	}

	router := gin.New()
	router.GET("/users/count", handler.Count)

	// Teste de contagem
	req := httptest.NewRequest("GET", "/users/count", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verificar se retornou a contagem correta
	var response map[string]int
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["count"] != 2 {
		t.Errorf("Expected count 2, got %d", response["count"])
	}

	// Verificar se os usuários foram criados corretamente
	count, err := mockUseCase.Count(context.Background())
	if err != nil {
		t.Fatalf("Failed to count users: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 users in mock, got %d", count)
	}
}

*/
