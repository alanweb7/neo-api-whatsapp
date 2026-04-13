package service

import (
	"context"
	"fmt"

	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/alan/baileys-saas/core-go/internal/infra/engineclient"
	"github.com/alan/baileys-saas/core-go/internal/repository"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SessionService struct {
	repo   *repository.SessionRepository
	engine *engineclient.Client
}

func NewSessionService(repo *repository.SessionRepository, engine *engineclient.Client) *SessionService {
	return &SessionService{repo: repo, engine: engine}
}

type EngineSessionResponse struct {
	SessionID string `json:"session_id"`
	Status    string `json:"status"`
	QRCode    string `json:"qr_code,omitempty"`
}

func (s *SessionService) Create(ctx context.Context, tenantID uuid.UUID, name string) (*domain.WhatsAppSession, error) {
	payload := map[string]any{"tenant_id": tenantID.String(), "name": name}
	var out EngineSessionResponse
	if err := s.engine.Post(ctx, "/internal/v1/sessions", payload, &out); err != nil {
		return nil, err
	}
	session := &domain.WhatsAppSession{
		TenantID:        tenantID,
		EngineSessionID: out.SessionID,
		Name:            name,
		Status:          domain.SessionStatus(out.Status),
		QRCode:          strptr(out.QRCode),
		Metadata:        datatypes.JSON([]byte(`{}`)),
	}
	if session.Status == "" {
		session.Status = domain.SessionCreated
	}
	if err := s.repo.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionService) Start(ctx context.Context, tenantID, id uuid.UUID) error {
	session, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	if err := s.engine.Post(ctx, fmt.Sprintf("/internal/v1/sessions/%s/start", session.EngineSessionID), map[string]any{}, nil); err != nil {
		return err
	}
	session.Status = domain.SessionStarting
	return s.repo.Update(ctx, session)
}

func (s *SessionService) Status(ctx context.Context, tenantID, id uuid.UUID) (map[string]any, error) {
	session, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	var out map[string]any
	if err := s.engine.Get(ctx, fmt.Sprintf("/internal/v1/sessions/%s/status", session.EngineSessionID), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *SessionService) GetQRCode(ctx context.Context, tenantID, id uuid.UUID) (map[string]any, error) {
	session, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	var out map[string]any
	if err := s.engine.Get(ctx, fmt.Sprintf("/internal/v1/sessions/%s/qr", session.EngineSessionID), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *SessionService) Reconnect(ctx context.Context, tenantID, id uuid.UUID) error {
	session, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	if err := s.engine.Post(ctx, fmt.Sprintf("/internal/v1/sessions/%s/reconnect", session.EngineSessionID), map[string]any{}, nil); err != nil {
		return err
	}
	session.Status = domain.SessionReconnecting
	return s.repo.Update(ctx, session)
}

func (s *SessionService) Disconnect(ctx context.Context, tenantID, id uuid.UUID) error {
	session, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	if err := s.engine.Post(ctx, fmt.Sprintf("/internal/v1/sessions/%s/disconnect", session.EngineSessionID), map[string]any{}, nil); err != nil {
		return err
	}
	session.Status = domain.SessionDisconnected
	return s.repo.Update(ctx, session)
}

func (s *SessionService) Remove(ctx context.Context, tenantID, id uuid.UUID) error {
	session, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	if err := s.engine.Post(ctx, fmt.Sprintf("/internal/v1/sessions/%s/remove", session.EngineSessionID), map[string]any{}, nil); err != nil {
		return err
	}
	return s.repo.Delete(ctx, tenantID, id)
}

func (s *SessionService) List(ctx context.Context, tenantID uuid.UUID) ([]domain.WhatsAppSession, error) {
	return s.repo.ListByTenant(ctx, tenantID)
}

func (s *SessionService) Get(ctx context.Context, tenantID, id uuid.UUID) (*domain.WhatsAppSession, error) {
	session, err := s.repo.GetByID(ctx, tenantID, id)
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	return session, err
}

func strptr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
