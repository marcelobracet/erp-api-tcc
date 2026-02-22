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

	for _, existing := range m.users {
		if string(existing.KeycloakID) == req.KeycloakID {
			return nil, userDomain.ErrUserAlreadyExists
		}
	}

	now := time.Now()
	u := &userDomain.User{
		ID:          dbtypes.UUID(req.KeycloakID),
		KeycloakID:  dbtypes.UUID(req.KeycloakID),
		TenantID:    dbtypes.UUID(req.TenantID),
		DisplayName: req.DisplayName,
		Email:       req.Email,


































































































































































































































































































}	}		t.Fatalf("expected count=2, got %d", payload["count"])	if payload["count"] != 2 {	}		t.Fatalf("failed to decode count: %v", err)	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {	var payload map[string]int	}		t.Fatalf("expected %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())	if w.Code != http.StatusOK {	router.ServeHTTP(w, req)	req = httptest.NewRequest("GET", "/users/count", nil)	w = httptest.NewRecorder()	}		t.Fatalf("expected 2 users, got %d", len(users))	if len(users) != 2 {	}		t.Fatalf("failed to decode list: %v", err)	if err := json.Unmarshal(w.Body.Bytes(), &users); err != nil {	var users []userDomain.User	}		t.Fatalf("expected %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())	if w.Code != http.StatusOK {	router.ServeHTTP(w, req)	req := httptest.NewRequest("GET", "/users?limit=10&offset=0", nil)	w := httptest.NewRecorder()	register()	register()	}		}			t.Fatalf("expected %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())		if w.Code != http.StatusCreated {		router.ServeHTTP(w, r)		r.Header.Set("Content-Type", "application/json")		r := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(b))		w := httptest.NewRecorder()		b, _ := json.Marshal(req)		}			Email:       strPtr("x@example.com"),			DisplayName: "User",			KeycloakID:  uuid.NewString(),			TenantID:    uuid.NewString(),		req := userDomain.CreateUserRequest{	register := func() {	router.GET("/users/count", h.Count)	router.GET("/users", h.List)	router.POST("/users/register", h.Register)	router := gin.New()	h := NewHandler(uc)	uc := NewMockUseCase()	gin.SetMode(gin.TestMode)func TestHandler_ListAndCount(t *testing.T) {}	}		t.Fatalf("expected %d, got %d: %s", http.StatusNotFound, w.Code, w.Body.String())	if w.Code != http.StatusNotFound {	router.ServeHTTP(w, req)	req := httptest.NewRequest("DELETE", "/users/missing", nil)	w := httptest.NewRecorder()	router.DELETE("/users/:id", h.Delete)	router := gin.New()	h := NewHandler(uc)	uc := NewMockUseCase()	gin.SetMode(gin.TestMode)func TestHandler_Delete(t *testing.T) {}	}		t.Fatalf("expected %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())	if w.Code != http.StatusOK {	router.ServeHTTP(w, req)	req.Header.Set("Content-Type", "application/json")	req = httptest.NewRequest("PUT", "/users/"+keycloakID, bytes.NewBuffer(body))	w = httptest.NewRecorder()	body, _ = json.Marshal(update)	update := userDomain.UpdateUserRequest{DisplayName: strPtr("New"), Email: &newEmail}	newEmail := "new@example.com"	}		t.Fatalf("expected %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())	if w.Code != http.StatusCreated {	router.ServeHTTP(w, req)	req.Header.Set("Content-Type", "application/json")	req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))	w := httptest.NewRecorder()	body, _ := json.Marshal(create)	}		Email:       strPtr("old@example.com"),		DisplayName: "Old",		KeycloakID:  keycloakID,		TenantID:    tenantID,	create := userDomain.CreateUserRequest{	tenantID := uuid.NewString()	keycloakID := uuid.NewString()	router.PUT("/users/:id", h.Update)	router.POST("/users/register", h.Register)	router := gin.New()	h := NewHandler(uc)	uc := NewMockUseCase()	gin.SetMode(gin.TestMode)func TestHandler_Update(t *testing.T) {}	}		t.Fatalf("expected email=%q, got %v", email, got.Email)	if got.Email == nil || *got.Email != email {	}		t.Fatalf("expected id=%q, got %q", keycloakID, string(got.ID))	if string(got.ID) != keycloakID {	}		t.Fatalf("failed to decode response: %v", err)	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {	var got userDomain.User	}		t.Fatalf("expected %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())	if w.Code != http.StatusOK {	router.ServeHTTP(w, req)	req := httptest.NewRequest("GET", "/users/profile", nil)	w := httptest.NewRecorder()	})		h.GetProfile(c)		c.Set("user_email", email)		c.Set("tenant_id", tenantID)		c.Set("user_id", keycloakID)	router.GET("/users/profile", func(c *gin.Context) {	router := gin.New()	email := "profile@example.com"	tenantID := uuid.NewString()	keycloakID := uuid.NewString()	h := NewHandler(uc)	uc := NewMockUseCase()	gin.SetMode(gin.TestMode)	t.Setenv("AUTH_PROVIDER", "keycloak")func TestHandler_GetProfile_KeycloakAutoProvision(t *testing.T) {}	}		t.Fatalf("expected %d, got %d: %s", http.StatusUnauthorized, w.Code, w.Body.String())	if w.Code != http.StatusUnauthorized {	router.ServeHTTP(w, req)	req := httptest.NewRequest("GET", "/users/profile", nil)	w := httptest.NewRecorder()	router.GET("/users/profile", h.GetProfile)	router := gin.New()	h := NewHandler(uc)	uc := NewMockUseCase()	gin.SetMode(gin.TestMode)func TestHandler_GetProfile_UnauthorizedWhenNoContext(t *testing.T) {}	}		t.Fatalf("expected %d, got %d: %s", http.StatusNotFound, w.Code, w.Body.String())	if w.Code != http.StatusNotFound {	router.ServeHTTP(w, req)	req := httptest.NewRequest("GET", "/users/does-not-exist", nil)	w := httptest.NewRecorder()	router.GET("/users/:id", h.GetByID)	router := gin.New()	h := NewHandler(uc)	uc := NewMockUseCase()	gin.SetMode(gin.TestMode)func TestHandler_GetByID_NotFound(t *testing.T) {}	}		t.Fatalf("expected %d, got %d: %s", http.StatusConflict, w.Code, w.Body.String())	if w.Code != http.StatusConflict {	router.ServeHTTP(w, req)	req.Header.Set("Content-Type", "application/json")	req = httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))	w = httptest.NewRecorder()	}		t.Fatalf("expected %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())	if w.Code != http.StatusCreated {	router.ServeHTTP(w, req)	w := httptest.NewRecorder()	req.Header.Set("Content-Type", "application/json")	req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))	body, _ := json.Marshal(reqBody)	}		Email:       &email,		DisplayName: "Test User",		KeycloakID:  keycloakID,		TenantID:    tenantID,	reqBody := userDomain.CreateUserRequest{	email := "test@example.com"	keycloakID := uuid.NewString()	tenantID := uuid.NewString()	router.POST("/users/register", h.Register)	router := gin.New()	h := NewHandler(uc)	uc := NewMockUseCase()	gin.SetMode(gin.TestMode)func TestHandler_Register_CreatesAndRejectsDuplicate(t *testing.T) {func strPtr(s string) *string { return &s }}	return len(m.users), nilfunc (m *MockUseCase) Count(ctx context.Context) (int, error) {}	return all[offset:end], nil	}		end = len(all)	if end > len(all) {	end := offset + limit	}		return []*userDomain.User{}, nil	if offset >= len(all) {	}		all = append(all, u)	for _, u := range m.users {	all := make([]*userDomain.User, 0, len(m.users))func (m *MockUseCase) List(ctx context.Context, limit, offset int) ([]*userDomain.User, error) {}	return nil	delete(m.users, id)	}		return userDomain.ErrUserNotFound	if _, ok := m.users[id]; !ok {func (m *MockUseCase) Delete(ctx context.Context, id string) error {}	return u, nil	u.UpdatedAt = time.Now()	}		u.Email = req.Email	if req.Email != nil {	}		u.DisplayName = *req.DisplayName	if req.DisplayName != nil {	}		return nil, userDomain.ErrUserNotFound	if !ok {	u, ok := m.users[id]func (m *MockUseCase) Update(ctx context.Context, id string, req *userDomain.UpdateUserRequest) (*userDomain.User, error) {}	return u, nil	}		return nil, userDomain.ErrUserNotFound	if !ok {	u, ok := m.users[id]func (m *MockUseCase) GetByID(ctx context.Context, id string) (*userDomain.User, error) {}	return u, nil	m.users[req.KeycloakID] = u	}		UpdatedAt:   now,		CreatedAt:   now,