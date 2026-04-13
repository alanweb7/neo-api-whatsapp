package service

import (
	"context"

	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/alan/baileys-saas/core-go/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantService struct{ repo *repository.TenantRepository }

func NewTenantService(repo *repository.TenantRepository) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) Create(ctx context.Context, t *domain.Tenant) error {
	return s.repo.Create(ctx, t)
}
func (s *TenantService) List(ctx context.Context) ([]domain.Tenant, error) { return s.repo.List(ctx) }
func (s *TenantService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	return t, err
}
func (s *TenantService) Update(ctx context.Context, t *domain.Tenant) error {
	return s.repo.Update(ctx, t)
}
