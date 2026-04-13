package repository

import (
	"context"
	"time"

	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type APIKeyRepository struct{ db *gorm.DB }

func NewAPIKeyRepository(db *gorm.DB) *APIKeyRepository { return &APIKeyRepository{db: db} }

func (r *APIKeyRepository) Create(ctx context.Context, apiKey *domain.ApiKey) error {
	return r.db.WithContext(ctx).Create(apiKey).Error
}

func (r *APIKeyRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]domain.ApiKey, error) {
	var keys []domain.ApiKey
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("created_at desc").Find(&keys).Error
	return keys, err
}

func (r *APIKeyRepository) Revoke(ctx context.Context, tenantID, keyID uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&domain.ApiKey{}).
		Where("id = ? and tenant_id = ?", keyID, tenantID).
		Updates(map[string]any{"revoked_at": &now}).Error
}
