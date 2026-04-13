package service

import (
	"context"

	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/alan/baileys-saas/core-go/internal/repository"
	"github.com/google/uuid"
)

type APIKeyService struct{ repo *repository.APIKeyRepository }

func NewAPIKeyService(repo *repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{repo: repo}
}

func (s *APIKeyService) Create(ctx context.Context, tenantID uuid.UUID, name string, createdBy *uuid.UUID) (*domain.ApiKey, string, error) {
	plain, prefix, hash, err := GenerateAPIKeyMaterial()
	if err != nil {
		return nil, "", err
	}
	k := &domain.ApiKey{TenantID: tenantID, Name: name, KeyPrefix: prefix, KeyHash: hash, CreatedByUID: createdBy}
	if err := s.repo.Create(ctx, k); err != nil {
		return nil, "", err
	}
	return k, plain, nil
}

func (s *APIKeyService) List(ctx context.Context, tenantID uuid.UUID) ([]domain.ApiKey, error) {
	return s.repo.ListByTenant(ctx, tenantID)
}

func (s *APIKeyService) Revoke(ctx context.Context, tenantID, keyID uuid.UUID) error {
	return s.repo.Revoke(ctx, tenantID, keyID)
}
