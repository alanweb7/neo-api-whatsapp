package repository

import (
	"context"

	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantRepository struct{ db *gorm.DB }

func NewTenantRepository(db *gorm.DB) *TenantRepository { return &TenantRepository{db: db} }

func (r *TenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *TenantRepository) List(ctx context.Context) ([]domain.Tenant, error) {
	var tenants []domain.Tenant
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&tenants).Error
	return tenants, err
}

func (r *TenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	var tenant domain.Tenant
	err := r.db.WithContext(ctx).First(&tenant, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	return r.db.WithContext(ctx).Save(tenant).Error
}
