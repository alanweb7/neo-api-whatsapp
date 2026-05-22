package middleware

import (
	"net/http"
	"strings"

	"github.com/alan/baileys-saas/core-go/internal/repository"
	"github.com/alan/baileys-saas/core-go/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha256"
)

func Auth(tokens *service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		claims, err := tokens.ParseAccess(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		userID, _ := uuid.Parse(claims.UserID)
		tenantID, _ := uuid.Parse(claims.TenantID)
		c.Set("user_id", userID)
		c.Set("tenant_id", tenantID)
		c.Set("auth_type", "jwt")
		c.Next()
	}
}

func AuthOrAPIKey(tokens *service.TokenService, apiKeyRepo *repository.APIKeyRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tentar API Key primeiro
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "" {
			// Validar API Key (sem validar expiração)
			apiKeyHash := hashAPIKey(apiKey)
			key, err := apiKeyRepo.GetByHash(c.Request.Context(), apiKeyHash)
			if err == nil && key != nil && key.RevokedAt == nil {
				c.Set("tenant_id", key.TenantID)
				c.Set("api_key_id", key.ID)
				c.Set("auth_type", "api_key")
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
			return
		}

		// Fallback para JWT (com validação de expiração)
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token or api key"})
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		claims, err := tokens.ParseAccess(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		userID, _ := uuid.Parse(claims.UserID)
		tenantID, _ := uuid.Parse(claims.TenantID)
		c.Set("user_id", userID)
		c.Set("tenant_id", tenantID)
		c.Set("auth_type", "jwt")
		c.Next()
	}
}

func hashAPIKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return string(h[:])
}
