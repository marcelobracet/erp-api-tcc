package middleware

import (
	"os"
	"strings"

	"erp-api/pkg/auth"

	"github.com/gin-gonic/gin"
)

type AuthProvider interface {
	Authenticate() gin.HandlerFunc
	RequireRole(requiredRole string) gin.HandlerFunc
	RequireAnyRole(requiredRoles ...string) gin.HandlerFunc
	OptionalAuth() gin.HandlerFunc
}

// NewAuthProvider seleciona o provedor de autenticação.
//
// Por padrão usa o JWT local. Para Keycloak, defina AUTH_PROVIDER=keycloak
// e configure KEYCLOAK_ISSUER e KEYCLOAK_AUDIENCE.
func NewAuthProvider(jwtManager *auth.JWTManager) AuthProvider {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("AUTH_PROVIDER"))) {
	case "keycloak":
		provider, err := NewKeycloakMiddlewareFromEnv()
		if err != nil {
			panic(err)
		}
		return provider
	default:
		return NewAuthMiddleware(jwtManager)
	}
}
