package auth

import (
	"testing"
	"time"
)

func TestJWTManager_GenerateAccessToken(t *testing.T) {
	secretKey := "test-secret-key"
	accessExpiry := 1 * time.Hour
	refreshExpiry := 7 * 24 * time.Hour

	jwtManager := NewJWTManager(secretKey, accessExpiry, refreshExpiry)

	userID := "user-123"
	email := "test@example.com"
	role := "user"
	provider := "local"

	token, err := jwtManager.GenerateAccessToken(userID, email, role, provider)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}
	if token == "" {
		t.Error("Expected token to be generated")
	}

	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}

	if claims.Role != role {
		t.Errorf("Expected Role %s, got %s", role, claims.Role)
	}
}

func TestJWTManager_GenerateRefreshToken(t *testing.T) {
	secretKey := "test-secret-key"
	accessExpiry := 1 * time.Hour
	refreshExpiry := 7 * 24 * time.Hour

	jwtManager := NewJWTManager(secretKey, accessExpiry, refreshExpiry)

	userID := "user-123"
	email := "test@example.com"
	role := "user"
	provider := "local"

	token, err := jwtManager.GenerateRefreshToken(userID, email, role, provider)
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}
	if token == "" {
		t.Error("Expected refresh token to be generated")
	}

	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}
}

func TestJWTManager_ValidateToken(t *testing.T) {
	secretKey := "test-secret-key"
	accessExpiry := 1 * time.Hour
	refreshExpiry := 7 * 24 * time.Hour

	jwtManager := NewJWTManager(secretKey, accessExpiry, refreshExpiry)

	userID := "user-123"
	email := "test@example.com"
	role := "user"
	provider := "test-provider"

	token, err := jwtManager.GenerateAccessToken(userID, email, role, provider)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}
	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}

	_, err = jwtManager.ValidateToken("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestJWTManager_RefreshAccessToken(t *testing.T) {
	secretKey := "test-secret-key"
	accessExpiry := 1 * time.Hour
	refreshExpiry := 7 * 24 * time.Hour

	jwtManager := NewJWTManager(secretKey, accessExpiry, refreshExpiry)

	userID := "user-123"
	email := "test@example.com"
	role := "user"

	deviceID := "test-device-id"
	refreshToken, err := jwtManager.GenerateRefreshToken(userID, email, role, deviceID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}
	newAccessToken, err := jwtManager.RefreshAccessToken(refreshToken)
	if err != nil {
		t.Fatalf("RefreshAccessToken() error = %v", err)
	}

	if newAccessToken == "" {
		t.Error("Expected new access token to be generated")
	}

	claims, err := jwtManager.ValidateToken(newAccessToken)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}
}

func TestJWTManager_GenerateTokenPair(t *testing.T) {
	secretKey := "test-secret-key"
	accessExpiry := 1 * time.Hour
	refreshExpiry := 7 * 24 * time.Hour

	jwtManager := NewJWTManager(secretKey, accessExpiry, refreshExpiry)

	userID := "user-123"
	email := "test@example.com"
	role := "user"

	deviceID := "test-device-id"
	tokenPair, err := jwtManager.GenerateTokenPair(userID, email, role, deviceID)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	if tokenPair.AccessToken == "" {
		t.Error("Expected access token to be generated")
	}

	if tokenPair.RefreshToken == "" {
		t.Error("Expected refresh token to be generated")
	}

	// Validar access token
	claims, err := jwtManager.ValidateToken(tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}

	// Validar refresh token
	refreshClaims, err := jwtManager.ValidateToken(tokenPair.RefreshToken)
	if err != nil {
		t.Fatalf("ValidateToken() refresh token error = %v", err)
	}

	if refreshClaims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, refreshClaims.UserID)
	}
}

func TestJWTManager_TokenExpiration(t *testing.T) {
	secretKey := "test-secret-key"
	accessExpiry := 1 * time.Millisecond // Token expira muito rápido
	refreshExpiry := 7 * 24 * time.Hour

	jwtManager := NewJWTManager(secretKey, accessExpiry, refreshExpiry)

	userID := "user-123"
	email := "test@example.com"
	role := "user"
	deviceID := "test-device-id"

	token, err := jwtManager.GenerateAccessToken(userID, email, role, deviceID)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}
	// Aguardar o token expirar
	time.Sleep(10 * time.Millisecond)

	// Tentar validar o token expirado
	_, err = jwtManager.ValidateToken(token)
	if err == nil {
		t.Error("Expected error for expired token")
	}
}
