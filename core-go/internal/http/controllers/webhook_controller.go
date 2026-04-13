package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/alan/baileys-saas/core-go/internal/http/dto"
	"github.com/alan/baileys-saas/core-go/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type WebhookController struct{ service *service.WebhookService }

func NewWebhookController(service *service.WebhookService) *WebhookController {
	return &WebhookController{service: service}
}

func (h *WebhookController) Create(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	var req dto.CreateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	jsonEvents := datatypes.JSON([]byte("[]"))
	if len(req.EventTypes) > 0 {
		if b, err := json.Marshal(req.EventTypes); err == nil {
			jsonEvents = datatypes.JSON(b)
		}
	}
	var secret *string
	if req.Secret != "" {
		secret = &req.Secret
	}
	hook := &domain.WebhookEndpoint{TenantID: tenantID, Name: req.Name, URL: req.URL, Secret: secret, IsActive: true, EventTypes: jsonEvents}
	if err := h.service.Create(c.Request.Context(), hook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, hook)
}

func (h *WebhookController) List(c *gin.Context) {
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

func (h *WebhookController) Update(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("webhookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req dto.UpdateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.Update(c.Request.Context(), tenantID, id, func(w *domain.WebhookEndpoint) {
		if req.Name != nil {
			w.Name = *req.Name
		}
		if req.URL != nil {
			w.URL = *req.URL
		}
		if req.Secret != nil {
			w.Secret = req.Secret
		}
		if req.IsActive != nil {
			w.IsActive = *req.IsActive
		}
		if len(req.EventTypes) > 0 {
			if b, e := json.Marshal(req.EventTypes); e == nil {
				w.EventTypes = datatypes.JSON(b)
			}
		}
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated": true})
}

func (h *WebhookController) Delete(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("webhookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.Delete(c.Request.Context(), tenantID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"removed": true})
}

func (h *WebhookController) ListDeliveries(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	items, err := h.service.ListDeliveries(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
