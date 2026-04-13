package controllers

import (
	"net/http"

	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/alan/baileys-saas/core-go/internal/http/dto"
	"github.com/alan/baileys-saas/core-go/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct{ service *service.UserService }

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

func (h *UserController) Create(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role := req.Role
	if role == "" {
		role = "member"
	}
	user := &domain.User{Email: req.Email, PasswordHash: req.Password, FullName: req.FullName, Status: "active"}
	if err := h.service.Create(c.Request.Context(), user, tenantID, role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *UserController) List(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	users, err := h.service.ListByTenant(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserController) Attach(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	var req dto.AttachUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, _ := uuid.Parse(req.UserID)
	role := req.Role
	if role == "" {
		role = "member"
	}
	if err := h.service.Attach(c.Request.Context(), tenantID, uid, role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"attached": true})
}
