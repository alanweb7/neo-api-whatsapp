package repository

import (
	"context"

	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]domain.User, error) {
	var users []domain.User
	err := r.db.WithContext(ctx).
		Table("users").
		Select("users.*").
		Joins("join tenant_users tu on tu.user_id = users.id").
		Where("tu.tenant_id = ?", tenantID).
		Order("users.created_at desc").
		Scan(&users).Error
	return users, err
}

func (r *UserRepository) AttachToTenant(ctx context.Context, tenantUser *domain.TenantUser) error {
	return r.db.WithContext(ctx).Create(tenantUser).Error
}
