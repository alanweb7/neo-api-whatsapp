package repository

import (
	"context"

	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WebhookRepository struct{ db *gorm.DB }

func NewWebhookRepository(db *gorm.DB) *WebhookRepository { return &WebhookRepository{db: db} }

func (r *WebhookRepository) Create(ctx context.Context, w *domain.WebhookEndpoint) error {
	return r.db.WithContext(ctx).Create(w).Error
}

func (r *WebhookRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]domain.WebhookEndpoint, error) {
	var hooks []domain.WebhookEndpoint
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("created_at desc").Find(&hooks).Error
	return hooks, err
}

func (r *WebhookRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.WebhookEndpoint, error) {
	var w domain.WebhookEndpoint
	err := r.db.WithContext(ctx).First(&w, "tenant_id = ? and id = ?", tenantID, id).Error
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WebhookRepository) Update(ctx context.Context, w *domain.WebhookEndpoint) error {
	return r.db.WithContext(ctx).Save(w).Error
}

func (r *WebhookRepository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("tenant_id = ? and id = ?", tenantID, id).Delete(&domain.WebhookEndpoint{}).Error
}

func (r *WebhookRepository) ListDeliveries(ctx context.Context, tenantID uuid.UUID) ([]domain.WebhookDelivery, error) {
	var deliveries []domain.WebhookDelivery
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("created_at desc").Limit(100).Find(&deliveries).Error
	return deliveries, err
}
