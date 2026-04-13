package repository

import (
	"context"

	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository struct{ db *gorm.DB }

func NewSessionRepository(db *gorm.DB) *SessionRepository { return &SessionRepository{db: db} }

func (r *SessionRepository) Create(ctx context.Context, session *domain.WhatsAppSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *SessionRepository) Update(ctx context.Context, session *domain.WhatsAppSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *SessionRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]domain.WhatsAppSession, error) {
	var sessions []domain.WhatsAppSession
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("created_at desc").Find(&sessions).Error
	return sessions, err
}

func (r *SessionRepository) GetByID(ctx context.Context, tenantID, sessionID uuid.UUID) (*domain.WhatsAppSession, error) {
	var s domain.WhatsAppSession
	err := r.db.WithContext(ctx).First(&s, "id = ? and tenant_id = ?", sessionID, tenantID).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepository) Delete(ctx context.Context, tenantID, sessionID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? and tenant_id = ?", sessionID, tenantID).Delete(&domain.WhatsAppSession{}).Error
}
