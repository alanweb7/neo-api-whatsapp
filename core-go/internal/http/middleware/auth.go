package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/alan/baileys-saas/core-go/internal/repository"
	"github.com/alan/baileys-saas/core-go/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		// Tentar API Key primeiro (aceita X-API-Key ou api-key)
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.GetHeader("api-key")
		}

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
	return hex.EncodeToString(h[:])
}

func InternalKey(internalAPIKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Internal-Key")
		if key == "" {
			key = c.GetHeader("api-key")
		}

		if key == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing internal api key"})
			return
		}

		if key != internalAPIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid internal api key"})
			return
		}

		c.Set("auth_type", "internal_key")
		c.Next()
	}
}

func EngineSessionAuth(sessionRepo *repository.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := c.Get("tenant_id")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "tenant_id not found in context"})
			return
		}

		engineSessionID := c.GetHeader("X-Engine-Session-ID")
		sessionIdParam := c.Param("sessionId")

		if engineSessionID == "" && sessionIdParam == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing session id"})
			return
		}

		var sessionID string
		if engineSessionID != "" {
			sessionID = engineSessionID
		} else {
			sessionID = sessionIdParam
		}

		parsedSessionID, err := uuid.Parse(sessionID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid session id format"})
			return
		}

		parsedTenantID, ok := tenantID.(uuid.UUID)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid tenant id type"})
			return
		}

		session, err := sessionRepo.GetByID(c.Request.Context(), parsedTenantID, parsedSessionID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "session not found"})
			return
		}

		c.Set("session_id", session.ID)
		c.Set("engine_session_id", session.EngineSessionID)
		c.Set("auth_type", "engine_session")
		c.Next()
	}
}
