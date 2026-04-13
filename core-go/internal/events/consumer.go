package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/alan/baileys-saas/core-go/internal/repository"
	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

type Consumer struct {
	redis    *goredis.Client
	sessions *repository.SessionRepository
	messages *repository.MessageRepository
	log      *logrus.Logger
}

func NewConsumer(redis *goredis.Client, sessions *repository.SessionRepository, messages *repository.MessageRepository, log *logrus.Logger) *Consumer {
	return &Consumer{redis: redis, sessions: sessions, messages: messages, log: log}
}

type Envelope struct {
	Version   string         `json:"version"`
	Type      string         `json:"type"`
	TenantID  string         `json:"tenant_id"`
	SessionID string         `json:"session_id,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
	Payload   map[string]any `json:"payload"`
}

func (c *Consumer) Start(ctx context.Context) {
	go c.consume(ctx)
}

func (c *Consumer) consume(ctx context.Context) {
	sub := c.redis.Subscribe(ctx, "wa.events.v1")
	defer sub.Close()
	ch := sub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			if msg == nil {
				continue
			}
			if err := c.handle(ctx, []byte(msg.Payload)); err != nil {
				c.log.WithError(err).Error("failed to process engine event")
			}
		}
	}
}

func (c *Consumer) handle(ctx context.Context, raw []byte) error {
	var e Envelope
	if err := json.Unmarshal(raw, &e); err != nil {
		return err
	}
	tenantID, err := uuid.Parse(e.TenantID)
	if err != nil {
		return err
	}

	switch e.Type {
	case "session.qr.updated", "session.connected", "session.disconnected", "session.failed", "session.reconnecting", "session.starting":
		return c.handleSessionEvent(ctx, tenantID, e)
	case "message.received", "message.sent", "message.delivery_update":
		return c.handleMessageEvent(ctx, tenantID, e)
	default:
		c.log.WithField("type", e.Type).Info("ignored event type")
		return nil
	}
}

func (c *Consumer) handleSessionEvent(ctx context.Context, tenantID uuid.UUID, e Envelope) error {
	if e.SessionID == "" {
		return fmt.Errorf("missing session_id")
	}
	for _, s := range mustListSessions(ctx, c.sessions, tenantID) {
		if s.EngineSessionID != e.SessionID {
			continue
		}
		clone := s
		if e.Type == "session.qr.updated" {
			if qr, ok := e.Payload["qr_code"].(string); ok {
				clone.QRCode = &qr
				clone.Status = domain.SessionQRPending
			}
		}
		if e.Type == "session.connected" {
			clone.Status = domain.SessionConnected
		}
		if e.Type == "session.disconnected" {
			clone.Status = domain.SessionDisconnected
		}
		if e.Type == "session.failed" {
			clone.Status = domain.SessionFailed
		}
		if e.Type == "session.reconnecting" {
			clone.Status = domain.SessionReconnecting
		}
		if e.Type == "session.starting" {
			clone.Status = domain.SessionStarting
		}
		return c.sessions.Update(ctx, &clone)
	}
	return nil
}

func (c *Consumer) handleMessageEvent(ctx context.Context, tenantID uuid.UUID, e Envelope) error {
	for _, s := range mustListSessions(ctx, c.sessions, tenantID) {
		if s.EngineSessionID != e.SessionID {
			continue
		}
		payloadBytes, _ := json.Marshal(e.Payload)
		msgType := "unknown"
		if v, ok := e.Payload["type"].(string); ok {
			msgType = v
		}
		direction := "inbound"
		if e.Type == "message.sent" {
			direction = "outbound"
		}
		entry := &domain.MessageLog{
			TenantID:          tenantID,
			WhatsAppSessionID: s.ID,
			Direction:         direction,
			MessageType:       msgType,
			Payload:           datatypes.JSON(payloadBytes),
		}
		return c.messages.CreateLog(ctx, entry)
	}
	return nil
}

func mustListSessions(ctx context.Context, repo *repository.SessionRepository, tenantID uuid.UUID) []domain.WhatsAppSession {
	rows, err := repo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil
	}
	return rows
}
