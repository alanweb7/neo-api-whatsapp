package service

import (
	"context"

	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/alan/baileys-saas/core-go/internal/repository"
	"github.com/google/uuid"
)

type WebhookService struct{ repo *repository.WebhookRepository }

func NewWebhookService(repo *repository.WebhookRepository) *WebhookService {
	return &WebhookService{repo: repo}
}

func (s *WebhookService) Create(ctx context.Context, hook *domain.WebhookEndpoint) error {
	return s.repo.Create(ctx, hook)
}
func (s *WebhookService) List(ctx context.Context, tenantID uuid.UUID) ([]domain.WebhookEndpoint, error) {
	return s.repo.ListByTenant(ctx, tenantID)
}
func (s *WebhookService) Update(ctx context.Context, tenantID, id uuid.UUID, updateFn func(*domain.WebhookEndpoint)) error {
	h, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	updateFn(h)
	return s.repo.Update(ctx, h)
}
func (s *WebhookService) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.Delete(ctx, tenantID, id)
}
func (s *WebhookService) ListDeliveries(ctx context.Context, tenantID uuid.UUID) ([]domain.WebhookDelivery, error) {
	return s.repo.ListDeliveries(ctx, tenantID)
}
