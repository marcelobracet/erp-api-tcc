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
	userUseCase "erp-api/internal/usecase/user"
	"erp-api/pkg/auth"

	"github.com/gin-gonic/gin"
)

// MockUseCase é um mock do usecase para testes
type MockUseCase struct {
	users map[string]*userDomain.User
}

// Garantir que MockUseCase implementa UseCaseInterface
var _ userUseCase.UseCaseInterface = (*MockUseCase)(nil)

func NewMockUseCase() *MockUseCase {
	return &MockUseCase{
		users: make(map[string]*userDomain.User),
	}
}

func (m *MockUseCase) Register(ctx context.Context, req *userDomain.CreateUserRequest) (*userDomain.User, error) {
	// Verificar se usuário já existe
	for _, user := range m.users {
		if user.Email == req.Email {
			return nil, userDomain.ErrUserAlreadyExists
		}
	}
	
	// Validar request
	if err := req.ValidateCreate(); err != nil {
		return nil, err
	}
	
	// Criar usuário
	user := &userDomain.User{
		ID:        "user-" + time.Now().Format("20060102150405") + "-" + req.Email,
		Email:     req.Email,
		Password:  "hashed_" + req.Password, // Simular hash
		Name:      req.Name,
		Role:      req.Role,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	m.users[user.ID] = user
	return user, nil
}

func (m *MockUseCase) Login(ctx context.Context, req *userDomain.LoginRequest) (*userDomain.LoginResponse, error) {
	// Validar request
	if err := req.ValidateLogin(); err != nil {
		return nil, err
	}
	
	// Buscar usuário por email
	var user *userDomain.User
	for _, u := range m.users {
		if u.Email == req.Email {
			user = u
			break
		}
	}
	
	if user == nil {
		return nil, userDomain.ErrInvalidCredentials
	}
	
	// Verificar senha (simulação simples)
	if user.Password != "hashed_"+req.Password {
		return nil, userDomain.ErrInvalidCredentials
	}
	
	// Verificar se usuário está ativo
	if !user.IsActive {
		return nil, userDomain.ErrUserInactive
	}
	
	// Simular tokens
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	tokenPair, err := jwtManager.GenerateTokenPair(user.ID, user.TenantID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	
	return &userDomain.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User:         *user,
	}, nil
}

func (m *MockUseCase) RefreshToken(ctx context.Context, req *userDomain.RefreshTokenRequest) (*userDomain.LoginResponse, error) {
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	
	// Validar refresh token
	claims, err := jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, userDomain.ErrInvalidToken
	}
	
	// Buscar usuário
	user, exists := m.users[claims.UserID]
	if !exists {
		return nil, userDomain.ErrUserNotFound
	}
	
	// Verificar se usuário está ativo
	if !user.IsActive {
		return nil, userDomain.ErrUserInactive
	}
	
	// Gerar novo access token
	accessToken, err := jwtManager.GenerateAccessToken(user.ID, user.TenantID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	
	return &userDomain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
		User:         *user,
	}, nil
}

func (m *MockUseCase) GetByID(ctx context.Context, id string) (*userDomain.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, userDomain.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUseCase) GetByEmail(ctx context.Context, email string) (*userDomain.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, userDomain.ErrUserNotFound
}

func (m *MockUseCase) Update(ctx context.Context, id string, req *userDomain.UpdateUserRequest) (*userDomain.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, userDomain.ErrUserNotFound
	}
	
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	
	user.UpdatedAt = time.Now()
	return user, nil
}

func (m *MockUseCase) Delete(ctx context.Context, id string) error {
	if _, exists := m.users[id]; !exists {
		return userDomain.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *MockUseCase) List(ctx context.Context, limit, offset int) ([]*userDomain.User, error) {
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

func (m *MockUseCase) Count(ctx context.Context) (int, error) {
	return len(m.users), nil
}

func (m *MockUseCase) ValidateToken(tokenString string) (*auth.Claims, error) {
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	return jwtManager.ValidateToken(tokenString)
}

func TestHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockUseCase := NewMockUseCase()
	handler := NewHandler(mockUseCase)
	
	router := gin.New()
	router.POST("/register", handler.Register)
	
	// Teste de registro válido
	reqBody := userDomain.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "user",
	}
	
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
	
	// Teste de registro com email duplicado
	req = httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, w.Code)
	}
}

func TestHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockUseCase := NewMockUseCase()
	handler := NewHandler(mockUseCase)
	
	// Criar usuário primeiro
	createReq := userDomain.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "user",
	}
	
	_, err := mockUseCase.Register(context.Background(), &createReq)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
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
		{Email: "user1@example.com", Password: "password1", Name: "User 1", Role: "user"},
		{Email: "user2@example.com", Password: "password2", Name: "User 2", Role: "user"},
		{Email: "user3@example.com", Password: "password3", Name: "User 3", Role: "user"},
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
		{Email: "user1@example.com", Password: "password1", Name: "User 1", Role: "user"},
		{Email: "user2@example.com", Password: "password2", Name: "User 2", Role: "user"},
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