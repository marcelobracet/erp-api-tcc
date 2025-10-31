package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"erp-api/pkg/auth"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware_Authenticate(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	middleware := NewAuthMiddleware(jwtManager)
	
	// Configurar Gin para teste
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Authorization header is required"}`,
		},
		{
			name:           "invalid authorization format",
			authHeader:     "InvalidFormat token123",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid authorization header format"}`,
		},
		{
			name:           "empty token",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Token is required"}`,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid or expired token"}`,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(middleware.Authenticate())
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})
			
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			
			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body %s, got %s", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestAuthMiddleware_Authenticate_ValidToken(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	middleware := NewAuthMiddleware(jwtManager)
	
	// Gerar token válido
	token, err := jwtManager.GenerateAccessToken("user-123", "test@example.com", "user", "test-device-id")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.Authenticate())
	router.GET("/test", func(c *gin.Context) {
		userID, exists := GetUserIDFromContext(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found"})
			return
		}
		
		userEmail, exists := GetUserEmailFromContext(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_email not found"})
			return
		}
		
		userRole, exists := GetUserRoleFromContext(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_role not found"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"user_id":    userID,
			"user_email": userEmail,
			"user_role":  userRole,
		})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	expectedBody := `{"user_email":"test@example.com","user_id":"user-123","user_role":"user"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestAuthMiddleware_RequireRole(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	middleware := NewAuthMiddleware(jwtManager)
	
	// Gerar token para usuário com role "user"
	token, err := jwtManager.GenerateAccessToken("user-123", "test@example.com", "user", "test-device-id")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequireRole("admin"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Deve retornar 403 Forbidden pois o usuário tem role "user" mas precisa de "admin"
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}
	
	expectedBody := `{"error":"Insufficient permissions"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestAuthMiddleware_RequireRole_Success(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	middleware := NewAuthMiddleware(jwtManager)
	
	// Gerar token para usuário com role "admin"
	token, err := jwtManager.GenerateAccessToken("user-123", "test@example.com", "admin", "test-device-id")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequireRole("admin"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Deve retornar 200 OK pois o usuário tem a role necessária
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthMiddleware_RequireAnyRole(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	middleware := NewAuthMiddleware(jwtManager)
	
	// Gerar token para usuário com role "user"
	token, err := jwtManager.GenerateAccessToken("user-123", "test@example.com", "user", "test-device-id")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequireAnyRole("admin", "manager"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Deve retornar 403 Forbidden pois o usuário tem role "user" mas precisa de "admin" ou "manager"
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestAuthMiddleware_RequireAnyRole_Success(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	middleware := NewAuthMiddleware(jwtManager)
	
	// Gerar token para usuário com role "manager"
	token, err := jwtManager.GenerateAccessToken("user-123", "test@example.com", "manager", "test-device-id")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequireAnyRole("admin", "manager"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Deve retornar 200 OK pois o usuário tem uma das roles necessárias
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthMiddleware_OptionalAuth(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour, 24*time.Hour)
	middleware := NewAuthMiddleware(jwtManager)
	
	// Gerar token válido
	token, err := jwtManager.GenerateAccessToken("user-123", "test@example.com", "user", "test-device-id")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.OptionalAuth())
	router.GET("/test", func(c *gin.Context) {
		userID, exists := GetUserIDFromContext(c)
		if exists {
			c.JSON(http.StatusOK, gin.H{
				"authenticated": true,
				"user_id":       userID,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"authenticated": false,
			})
		}
	})
	
	// Teste com token válido
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	expectedBody := `{"authenticated":true,"user_id":"user-123"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
	
	// Teste sem token
	req = httptest.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	expectedBody = `{"authenticated":false}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	
	// Teste sem user_id no contexto
	userID, exists := GetUserIDFromContext(c)
	if exists {
		t.Error("Expected user_id to not exist")
	}
	
	// Adicionar user_id ao contexto
	c.Set("user_id", "test-user-id")
	
	userID, exists = GetUserIDFromContext(c)
	if !exists {
		t.Error("Expected user_id to exist")
	}
	
	if userID != "test-user-id" {
		t.Errorf("Expected user_id 'test-user-id', got '%s'", userID)
	}
}

func TestGetUserEmailFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	
	// Teste sem user_email no contexto
	userEmail, exists := GetUserEmailFromContext(c)
	if exists {
		t.Error("Expected user_email to not exist")
	}
	
	// Adicionar user_email ao contexto
	c.Set("user_email", "test@example.com")
	
	userEmail, exists = GetUserEmailFromContext(c)
	if !exists {
		t.Error("Expected user_email to exist")
	}
	
	if userEmail != "test@example.com" {
		t.Errorf("Expected user_email 'test@example.com', got '%s'", userEmail)
	}
}

func TestGetUserRoleFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	
	// Teste sem user_role no contexto
	userRole, exists := GetUserRoleFromContext(c)
	if exists {
		t.Error("Expected user_role to not exist")
	}
	
	// Adicionar user_role ao contexto
	c.Set("user_role", "admin")
	
	userRole, exists = GetUserRoleFromContext(c)
	if !exists {
		t.Error("Expected user_role to exist")
	}
	
	if userRole != "admin" {
		t.Errorf("Expected user_role 'admin', got '%s'", userRole)
	}
} 