package controllers

import (
	"net/http"

	"github.com/alan/baileys-saas/core-go/internal/http/dto"
	"github.com/alan/baileys-saas/core-go/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessageController struct{ service *service.MessageService }

func NewMessageController(service *service.MessageService) *MessageController {
	return &MessageController{service: service}
}

func (h *MessageController) SendText(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	var req dto.SendTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sid, _ := uuid.Parse(req.SessionID)
	out, err := h.service.SendText(c.Request.Context(), tenantID, sid, req.To, req.Text)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *MessageController) SendImage(c *gin.Context)    { h.sendMedia(c, "image") }
func (h *MessageController) SendDocument(c *gin.Context) { h.sendMedia(c, "document") }
func (h *MessageController) SendAudio(c *gin.Context)    { h.sendMedia(c, "audio") }

func (h *MessageController) sendMedia(c *gin.Context, t string) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	var req dto.SendMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sid, _ := uuid.Parse(req.SessionID)
	payload := map[string]any{"to": req.To, "media_url": req.MediaURL, "caption": req.Caption, "file_name": req.FileName}
	out, err := h.service.SendMedia(c.Request.Context(), tenantID, sid, t, payload)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *MessageController) ListLogs(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	items, err := h.service.ListLogs(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
