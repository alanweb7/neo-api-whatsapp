package repository

import (
	"context"

	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository struct{ db *gorm.DB }

func NewMessageRepository(db *gorm.DB) *MessageRepository { return &MessageRepository{db: db} }

func (r *MessageRepository) CreateLog(ctx context.Context, log *domain.MessageLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *MessageRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]domain.MessageLog, error) {
	var logs []domain.MessageLog
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("created_at desc").Limit(200).Find(&logs).Error
	return logs, err
}
