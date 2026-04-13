package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/alan/baileys-saas/core-go/internal/infra/engineclient"
	"github.com/alan/baileys-saas/core-go/internal/repository"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type MessageService struct {
	repo     *repository.MessageRepository
	sessions *repository.SessionRepository
	engine   *engineclient.Client
}

func NewMessageService(repo *repository.MessageRepository, sessions *repository.SessionRepository, engine *engineclient.Client) *MessageService {
	return &MessageService{repo: repo, sessions: sessions, engine: engine}
}

func (s *MessageService) SendText(ctx context.Context, tenantID, sessionID uuid.UUID, to, text string) (map[string]any, error) {
	session, err := s.sessions.GetByID(ctx, tenantID, sessionID)
	if err != nil {
		return nil, err
	}
	payload := map[string]any{"to": to, "text": text}
	var out map[string]any
	if err := s.engine.Post(ctx, fmt.Sprintf("/internal/v1/sessions/%s/messages/text", session.EngineSessionID), payload, &out); err != nil {
		return nil, err
	}
	messageType := "text"
	direction := "outbound"
	log := &domain.MessageLog{TenantID: tenantID, WhatsAppSessionID: session.ID, Direction: direction, MessageType: messageType, ToNumber: &to, Payload: datatypes.JSON(mustJSON(payload))}
	_ = s.repo.CreateLog(ctx, log)
	return out, nil
}

func (s *MessageService) SendMedia(ctx context.Context, tenantID, sessionID uuid.UUID, msgType string, payload map[string]any) (map[string]any, error) {
	session, err := s.sessions.GetByID(ctx, tenantID, sessionID)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := s.engine.Post(ctx, fmt.Sprintf("/internal/v1/sessions/%s/messages/%s", session.EngineSessionID, msgType), payload, &out); err != nil {
		return nil, err
	}
	log := &domain.MessageLog{TenantID: tenantID, WhatsAppSessionID: session.ID, Direction: "outbound", MessageType: msgType, Payload: datatypes.JSON(mustJSON(payload))}
	_ = s.repo.CreateLog(ctx, log)
	return out, nil
}

func (s *MessageService) SendButtons(ctx context.Context, tenantID, sessionID uuid.UUID, payload map[string]any) (map[string]any, error) {
	session, err := s.sessions.GetByID(ctx, tenantID, sessionID)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := s.engine.Post(ctx, fmt.Sprintf("/internal/v1/sessions/%s/messages/buttons", session.EngineSessionID), payload, &out); err != nil {
		return nil, err
	}
	log := &domain.MessageLog{TenantID: tenantID, WhatsAppSessionID: session.ID, Direction: "outbound", MessageType: "buttons", Payload: datatypes.JSON(mustJSON(payload))}
	_ = s.repo.CreateLog(ctx, log)
	return out, nil
}

func (s *MessageService) ListLogs(ctx context.Context, tenantID uuid.UUID) ([]domain.MessageLog, error) {
	return s.repo.ListByTenant(ctx, tenantID)
}

func mustJSON(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		return []byte("{}")
	}
	return data
}
