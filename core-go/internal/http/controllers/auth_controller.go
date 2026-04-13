package controllers

import (
	"net/http"

	"github.com/alan/baileys-saas/core-go/internal/http/dto"
	"github.com/alan/baileys-saas/core-go/internal/repository"
	"github.com/alan/baileys-saas/core-go/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController struct {
	auth  *service.AuthService
	users *repository.UserRepository
}

func NewAuthController(auth *service.AuthService, users *repository.UserRepository) *AuthController {
	return &AuthController{auth: auth, users: users}
}

func (h *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID, _ := uuid.Parse(req.TenantID)
	access, refresh, err := h.auth.Login(c.Request.Context(), req.Email, req.Password, tenantID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": refresh, "token_type": "Bearer"})
}

func (h *AuthController) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID, _ := uuid.Parse(req.TenantID)
	access, refresh, err := h.auth.Refresh(c.Request.Context(), req.RefreshToken, tenantID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": refresh, "token_type": "Bearer"})
}

func (h *AuthController) Me(c *gin.Context) {
	uid, ok := userIDFromCtx(c)
	if !ok {
		return
	}
	user, err := h.users.GetByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
