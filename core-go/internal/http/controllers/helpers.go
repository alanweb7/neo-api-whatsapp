package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func tenantIDFromCtx(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get("tenant_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found in token"})
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	if !ok || id == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid tenant in token"})
		return uuid.Nil, false
	}
	return id, true
}

func userIDFromCtx(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in token"})
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	if !ok || id == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user in token"})
		return uuid.Nil, false
	}
	return id, true
}
