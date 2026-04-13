package service

import (
	"context"

	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/alan/baileys-saas/core-go/internal/repository"
	"github.com/google/uuid"
)

type UserService struct{ repo *repository.UserRepository }

func NewUserService(repo *repository.UserRepository) *UserService { return &UserService{repo: repo} }

func (s *UserService) Create(ctx context.Context, user *domain.User, tenantID uuid.UUID, role string) error {
	h, err := HashPassword(user.PasswordHash)
	if err != nil {
		return err
	}
	user.PasswordHash = h
	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}
	return s.repo.AttachToTenant(ctx, &domain.TenantUser{TenantID: tenantID, UserID: user.ID, Role: role})
}

func (s *UserService) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]domain.User, error) {
	return s.repo.ListByTenant(ctx, tenantID)
}

func (s *UserService) Attach(ctx context.Context, tenantID, userID uuid.UUID, role string) error {
	return s.repo.AttachToTenant(ctx, &domain.TenantUser{TenantID: tenantID, UserID: userID, Role: role})
}
