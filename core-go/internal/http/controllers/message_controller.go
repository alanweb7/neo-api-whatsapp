package controllers

import (
	"net/http"
	"strings"

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
func (h *MessageController) SendCarousel(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	var req dto.SendCarouselRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	target := req.JID
	if target == "" {
		target = req.To
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "jid or to is required"})
		return
	}

	for cardIdx, card := range req.Cards {
		hasQuickReply := false
		hasCTA := false
		for btnIdx, button := range card.Buttons {
			switch button.Type {
			case "quick_reply":
				hasQuickReply = true
			case "cta_url", "cta_call", "cta_copy":
				hasCTA = true
			}
			if button.Type == "cta_url" && strings.TrimSpace(button.URL) == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "url is required for cta_url buttons", "card": cardIdx, "button": btnIdx})
				return
			}
			if button.Type == "cta_call" && strings.TrimSpace(button.PhoneNumber) == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "phoneNumber is required for cta_call buttons", "card": cardIdx, "button": btnIdx})
				return
			}
			if button.Type == "cta_copy" && strings.TrimSpace(button.CopyCode) == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "copyCode is required for cta_copy buttons", "card": cardIdx, "button": btnIdx})
				return
			}
		}
		if hasQuickReply && hasCTA {
			c.JSON(http.StatusBadRequest, gin.H{"error": "do not mix quick_reply with cta button types in the same card", "card": cardIdx})
			return
		}
	}

	sid, _ := uuid.Parse(req.SessionID)
	payload := map[string]any{
		"jid":           target,
		"text":          req.Text,
		"footer":        req.Footer,
		"cards":         req.Cards,
		"fallback_text": req.FallbackText,
	}
	out, err := h.service.SendCarousel(c.Request.Context(), tenantID, sid, payload)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *MessageController) SendButtons(c *gin.Context) {
	tenantID, ok := tenantIDFromCtx(c)
	if !ok {
		return
	}
	var req dto.SendButtonsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	target := req.JID
	if target == "" {
		target = req.To
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "jid or to is required"})
		return
	}
	hasQuickReply := false
	hasCTA := false
	for _, button := range req.Buttons {
		switch button.Type {
		case "quick_reply":
			hasQuickReply = true
		case "cta_url", "cta_call", "cta_copy":
			hasCTA = true
		}
		if button.Type == "cta_url" && strings.TrimSpace(button.URL) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "url is required for cta_url buttons"})
			return
		}
		if button.Type == "cta_call" && strings.TrimSpace(button.PhoneNumber) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "phoneNumber is required for cta_call buttons"})
			return
		}
		if button.Type == "cta_copy" && strings.TrimSpace(button.CopyCode) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "copyCode is required for cta_copy buttons"})
			return
		}
	}
	if hasQuickReply && hasCTA {
		c.JSON(http.StatusBadRequest, gin.H{"error": "do not mix quick_reply with cta button types in the same payload"})
		return
	}
	if hasCTA && len(req.Buttons) > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cta payload supports at most 3 buttons"})
		return
	}
	if hasQuickReply && len(req.Buttons) > 16 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quick_reply payload supports at most 16 buttons"})
		return
	}
	sid, _ := uuid.Parse(req.SessionID)
	payload := map[string]any{
		"jid":           target,
		"text":          req.Text,
		"footer":        req.Footer,
		"buttons":       req.Buttons,
		"fallback_text": req.FallbackText,
	}
	out, err := h.service.SendButtons(c.Request.Context(), tenantID, sid, payload)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

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
