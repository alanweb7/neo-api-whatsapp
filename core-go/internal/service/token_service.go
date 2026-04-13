package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewTokenService(accessSecret, refreshSecret string, accessTTLMin, refreshTTLDays int) *TokenService {
	return &TokenService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     time.Duration(accessTTLMin) * time.Minute,
		refreshTTL:    time.Duration(refreshTTLDays) * 24 * time.Hour,
	}
}

type AccessClaims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *TokenService) GeneratePair(userID, tenantID uuid.UUID) (string, string, error) {
	now := time.Now().UTC()
	accessClaims := AccessClaims{
		UserID:   userID.String(),
		TenantID: tenantID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
			Subject:   userID.String(),
		},
	}
	refreshClaims := RefreshClaims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTTL)),
			Subject:   userID.String(),
		},
	}

	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(s.accessSecret)
	if err != nil {
		return "", "", err
	}
	refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(s.refreshSecret)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func (s *TokenService) ParseAccess(tokenString string) (*AccessClaims, error) {
	parsed, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(_ *jwt.Token) (interface{}, error) {
		return s.accessSecret, nil
	})
	if err != nil || !parsed.Valid {
		return nil, ErrUnauthorized
	}
	claims, ok := parsed.Claims.(*AccessClaims)
	if !ok {
		return nil, ErrUnauthorized
	}
	return claims, nil
}

func (s *TokenService) ParseRefresh(tokenString string) (*RefreshClaims, error) {
	parsed, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(_ *jwt.Token) (interface{}, error) {
		return s.refreshSecret, nil
	})
	if err != nil || !parsed.Valid {
		return nil, ErrUnauthorized
	}
	claims, ok := parsed.Claims.(*RefreshClaims)
	if !ok {
		return nil, ErrUnauthorized
	}
	return claims, nil
}

func GenerateAPIKeyMaterial() (plain string, prefix string, hash string, err error) {
	raw := make([]byte, 32)
	if _, err = rand.Read(raw); err != nil {
		return "", "", "", err
	}
	enc := base64.RawURLEncoding.EncodeToString(raw)
	plain = fmt.Sprintf("wak_%s", enc)
	if len(plain) >= 12 {
		prefix = plain[:12]
	} else {
		prefix = plain
	}
	sum := sha256.Sum256([]byte(plain))
	hash = hex.EncodeToString(sum[:])
	return plain, prefix, hash, nil
}
