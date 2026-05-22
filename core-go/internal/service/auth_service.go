package service

import (
	"context"

	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/alan/baileys-saas/core-go/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	users   *repository.UserRepository
	tenants *repository.TenantRepository
	tokens  *TokenService
}

func NewAuthService(users *repository.UserRepository, tenants *repository.TenantRepository, tokens *TokenService) *AuthService {
	return &AuthService{users: users, tenants: tenants, tokens: tokens}
}

func (s *AuthService) Login(ctx context.Context, email, password string, tenantID uuid.UUID) (string, string, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", "", ErrUnauthorized
		}
		return "", "", err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", "", ErrUnauthorized
	}
	access, refresh, err := s.tokens.GeneratePair(user.ID, tenantID)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string, tenantID uuid.UUID) (string, string, error) {
	claims, err := s.tokens.ParseRefresh(refreshToken)
	if err != nil {
		return "", "", err
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return "", "", ErrUnauthorized
	}
	if _, err := s.users.GetByID(ctx, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", "", ErrUnauthorized
		}
		return "", "", err
	}
	return s.tokens.GeneratePair(userID, tenantID)
}

func (s *AuthService) Register(ctx context.Context, email, password, fullName, tenantName string) (*domain.User, uuid.UUID, error) {
	// Verificar se email já existe
	_, err := s.users.GetByEmail(ctx, email)
	if err == nil {
		return nil, uuid.Nil, ErrUnauthorized
	}
	if err != gorm.ErrRecordNotFound {
		return nil, uuid.Nil, err
	}

	// Hash da senha
	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, uuid.Nil, err
	}

	// Criar usuário
	user := &domain.User{
		BaseModel: domain.BaseModel{ID: uuid.New()},
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		Status:       "active",
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, uuid.Nil, err
	}

	// Criar tenant padrão
	tenantSlug := email[:len(email)-len("@")] + "-tenant"
	tenant := &domain.Tenant{
		BaseModel: domain.BaseModel{ID: uuid.New()},
		Name:   tenantName,
		Slug:   tenantSlug,
		Status: domain.TenantStatusActive,
	}
	if err := s.tenants.Create(ctx, tenant); err != nil {
		return nil, uuid.Nil, err
	}

	return user, tenant.ID, nil
}

func HashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}
