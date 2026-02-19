package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

type KeycloakMiddleware struct {
	issuer    string
	jwksURL   string
	audiences []string

	once     sync.Once
	verifier *oidc.IDTokenVerifier
	initErr  error
}

func NewKeycloakMiddlewareFromEnv() (*KeycloakMiddleware, error) {
	issuer := strings.TrimSpace(os.Getenv("KEYCLOAK_ISSUER"))
	if issuer == "" {
		return nil, errors.New("KEYCLOAK_ISSUER is required when AUTH_PROVIDER=keycloak")
	}

	jwksURL := strings.TrimSpace(os.Getenv("KEYCLOAK_JWKS_URL"))

	audienceRaw := strings.TrimSpace(os.Getenv("KEYCLOAK_AUDIENCE"))
	var audiences []string
	if audienceRaw != "" {
		for _, a := range strings.Split(audienceRaw, ",") {
			value := strings.TrimSpace(a)
			if value != "" {
				audiences = append(audiences, value)
			}
		}
	}

	return &KeycloakMiddleware{issuer: issuer, jwksURL: jwksURL, audiences: audiences}, nil
}

func (m *KeycloakMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := m.validateBearerToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", claims.Subject)
		if claims.Email != "" {
			c.Set("user_email", claims.Email)
		}
		if claims.TenantID != "" {
			c.Set("tenant_id", claims.TenantID)
		}

		// Mantém compatibilidade com o modelo atual (string "user_role")
		// e também salva a lista completa ("user_roles").
		if len(claims.Roles) > 0 {
			c.Set("user_role", claims.Roles[0])
			c.Set("user_roles", claims.Roles)
		}

		c.Next()
	}
}

func (m *KeycloakMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := m.validateBearerToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !containsString(claims.Roles, requiredRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.Subject)
		c.Set("tenant_id", claims.TenantID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", requiredRole)
		c.Set("user_roles", claims.Roles)
		c.Next()
	}
}

func (m *KeycloakMiddleware) RequireAnyRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := m.validateBearerToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		has := false
		for _, role := range requiredRoles {
			if containsString(claims.Roles, role) {
				has = true
				c.Set("user_role", role)
				break
			}
		}
		if !has {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.Subject)
		c.Set("tenant_id", claims.TenantID)
		c.Set("user_email", claims.Email)
		c.Set("user_roles", claims.Roles)
		c.Next()
	}
}

func (m *KeycloakMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, _ = m.trySetContextFromBearerToken(c)
		c.Next()
	}
}

type keycloakClaims struct {
	Subject  string
	Email    string
	TenantID string
	Roles    []string
}

func (m *KeycloakMiddleware) ensureVerifier() error {
	m.once.Do(func() {
		httpClient := &http.Client{Timeout: 5 * time.Second}
		ctx := oidc.ClientContext(context.Background(), httpClient)

		cfg := &oidc.Config{SkipClientIDCheck: true}
		if m.jwksURL != "" {
			keySet := oidc.NewRemoteKeySet(ctx, m.jwksURL)
			m.verifier = oidc.NewVerifier(m.issuer, keySet, cfg)
			return
		}

		provider, err := oidc.NewProvider(ctx, m.issuer)
		if err != nil {
			m.initErr = err
			return
		}
		m.verifier = provider.Verifier(cfg)
	})

	return m.initErr
}

func (m *KeycloakMiddleware) validateBearerToken(c *gin.Context) (*keycloakClaims, error) {
	if err := m.ensureVerifier(); err != nil {
		return nil, err
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("Authorization header is required")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if token == "" {
		return nil, errors.New("Token is required")
	}

	idToken, err := m.verifier.Verify(c.Request.Context(), token)
	if err != nil {
		return nil, errors.New("Invalid or expired token")
	}

	var raw map[string]any
	if err := idToken.Claims(&raw); err != nil {
		return nil, err
	}

	claims := &keycloakClaims{
		Subject: idToken.Subject,
		Email:   getString(raw, "email"),
		Roles:   extractRoles(raw, os.Getenv("KEYCLOAK_CLIENT_ID")),
	}

	// Opcional: aud check
	if len(m.audiences) > 0 && !audienceOrAzpAllowed(raw, m.audiences) {
		return nil, errors.New("Invalid token audience")
	}

	claims.TenantID = getString(raw, "tenant_id")
	if claims.TenantID == "" {
		// Para multi-tenant, este claim é obrigatório.
		return nil, errors.New("tenant_id not found in token")
	}

	return claims, nil
}

func (m *KeycloakMiddleware) trySetContextFromBearerToken(c *gin.Context) (*keycloakClaims, error) {
	claims, err := m.validateBearerToken(c)
	if err != nil {
		return nil, err
	}
	c.Set("user_id", claims.Subject)
	c.Set("tenant_id", claims.TenantID)
	c.Set("user_email", claims.Email)
	if len(claims.Roles) > 0 {
		c.Set("user_role", claims.Roles[0])
		c.Set("user_roles", claims.Roles)
	}
	return claims, nil
}

func extractRoles(raw map[string]any, clientID string) []string {
	// Realm roles
	var roles []string
	if realmAccess, ok := raw["realm_access"].(map[string]any); ok {
		roles = append(roles, asStringSlice(realmAccess["roles"])...)
	}

	// Client roles: resource_access[clientID].roles
	if clientID != "" {
		if ra, ok := raw["resource_access"].(map[string]any); ok {
			if client, ok := ra[clientID].(map[string]any); ok {
				roles = append(roles, asStringSlice(client["roles"])...)
			}
		}
	}

	return uniqueStrings(roles)
}

func getString(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	var out string
	_ = json.Unmarshal(b, &out)
	return out
}

func asStringSlice(v any) []string {
	if v == nil {
		return nil
	}
	if ss, ok := v.([]string); ok {
		return ss
	}
	if a, ok := v.([]any); ok {
		out := make([]string, 0, len(a))
		for _, item := range a {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

func audienceAllowed(aud any, allowed []string) bool {
	if aud == nil {
		return false
	}
	if s, ok := aud.(string); ok {
		return containsString(allowed, s)
	}
	if arr, ok := aud.([]any); ok {
		for _, item := range arr {
			if s, ok := item.(string); ok {
				if containsString(allowed, s) {
					return true
				}
			}
		}
		return false
	}
	return false
}

func audienceOrAzpAllowed(raw map[string]any, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	if audienceAllowed(raw["aud"], allowed) {
		return true
	}
	azp := getString(raw, "azp")
	if azp != "" && containsString(allowed, azp) {
		return true
	}
	return false
}

func containsString(list []string, v string) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}
	return false
}

func uniqueStrings(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, s := range in {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
