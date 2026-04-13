package controllers

import (
	"net/http"

	"github.com/alan/baileys-saas/core-go/internal/http/dto"
	"github.com/alan/baileys-saas/core-go/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SessionController struct{ service *service.SessionService }

func NewSessionController(service *service.SessionService) *SessionController {
	return &SessionController{service: service}
}

func (h *SessionController) Create(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	var req dto.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	session, err := h.service.Create(c.Request.Context(), tenantID, req.Name)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, session)
}

func (h *SessionController) Start(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	sid, err := uuid.Parse(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	if err := h.service.Start(c.Request.Context(), tenantID, sid); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"started": true})
}

func (h *SessionController) List(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	items, err := h.service.List(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *SessionController) Get(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	sid, err := uuid.Parse(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	item, err := h.service.Get(c.Request.Context(), tenantID, sid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *SessionController) QR(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	sid, err := uuid.Parse(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	out, err := h.service.GetQRCode(c.Request.Context(), tenantID, sid)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *SessionController) Status(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	sid, err := uuid.Parse(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	out, err := h.service.Status(c.Request.Context(), tenantID, sid)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *SessionController) Reconnect(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	sid, err := uuid.Parse(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	if err := h.service.Reconnect(c.Request.Context(), tenantID, sid); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reconnecting": true})
}

func (h *SessionController) Disconnect(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	sid, err := uuid.Parse(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	if err := h.service.Disconnect(c.Request.Context(), tenantID, sid); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"disconnected": true})
}

func (h *SessionController) Remove(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	sid, err := uuid.Parse(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	if err := h.service.Remove(c.Request.Context(), tenantID, sid); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"removed": true})
}
