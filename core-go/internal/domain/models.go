package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TenantStatus string
type SessionStatus string

const (
	TenantStatusActive   TenantStatus = "active"
	TenantStatusInactive TenantStatus = "inactive"

	SessionCreated      SessionStatus = "created"
	SessionStarting     SessionStatus = "starting"
	SessionQRPending    SessionStatus = "qr_pending"
	SessionConnected    SessionStatus = "connected"
	SessionDisconnected SessionStatus = "disconnected"
	SessionFailed       SessionStatus = "failed"
	SessionReconnecting SessionStatus = "reconnecting"
)

type BaseModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

type Tenant struct {
	BaseModel
	Name   string         `gorm:"size:120;not null" json:"name"`
	Slug   string         `gorm:"size:80;uniqueIndex;not null" json:"slug"`
	Status TenantStatus   `gorm:"size:20;not null;default:active" json:"status"`
	PlanID *uuid.UUID     `gorm:"type:uuid" json:"plan_id,omitempty"`
	Config datatypes.JSON `gorm:"type:jsonb" json:"config"`
}

type Plan struct {
	BaseModel
	Code      string         `gorm:"size:40;uniqueIndex;not null" json:"code"`
	Name      string         `gorm:"size:80;not null" json:"name"`
	Limits    datatypes.JSON `gorm:"type:jsonb" json:"limits"`
	IsDefault bool           `gorm:"not null;default:false" json:"is_default"`
}

type User struct {
	BaseModel
	Email        string `gorm:"size:180;uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	FullName     string `gorm:"size:120;not null" json:"full_name"`
	Status       string `gorm:"size:20;not null;default:active" json:"status"`
}

type TenantUser struct {
	BaseModel
	TenantID uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"`
	UserID   uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	Role     string    `gorm:"size:30;not null;default:member" json:"role"`
}

type ApiKey struct {
	BaseModel
	TenantID     uuid.UUID  `gorm:"type:uuid;index;not null" json:"tenant_id"`
	Name         string     `gorm:"size:80;not null" json:"name"`
	KeyPrefix    string     `gorm:"size:20;index;not null" json:"key_prefix"`
	KeyHash      string     `gorm:"size:255;not null" json:"-"`
	LastUsedAt   *time.Time `json:"last_used_at,omitempty"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
	CreatedByUID *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`
}

type WhatsAppSession struct {
	BaseModel
	TenantID        uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	EngineSessionID string         `gorm:"size:100;index;not null" json:"engine_session_id"`
	Name            string         `gorm:"size:120;not null" json:"name"`
	Status          SessionStatus  `gorm:"size:30;index;not null;default:created" json:"status"`
	Phone           *string        `gorm:"size:40" json:"phone,omitempty"`
	PushName        *string        `gorm:"size:120" json:"push_name,omitempty"`
	QRCode          *string        `gorm:"type:text" json:"qr_code,omitempty"`
	LastSeenAt      *time.Time     `json:"last_seen_at,omitempty"`
	ConnectedAt     *time.Time     `json:"connected_at,omitempty"`
	DisconnectedAt  *time.Time     `json:"disconnected_at,omitempty"`
	FailureReason   *string        `gorm:"type:text" json:"failure_reason,omitempty"`
	Metadata        datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
}

type WebhookEndpoint struct {
	BaseModel
	TenantID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	Name       string         `gorm:"size:80;not null" json:"name"`
	URL        string         `gorm:"size:300;not null" json:"url"`
	Secret     *string        `gorm:"size:120" json:"-"`
	IsActive   bool           `gorm:"not null;default:true" json:"is_active"`
	EventTypes datatypes.JSON `gorm:"type:jsonb" json:"event_types"`
}

type WebhookDelivery struct {
	BaseModel
	TenantID          uuid.UUID  `gorm:"type:uuid;index;not null" json:"tenant_id"`
	WebhookEndpointID uuid.UUID  `gorm:"type:uuid;index;not null" json:"webhook_endpoint_id"`
	EventType         string     `gorm:"size:80;index;not null" json:"event_type"`
	Payload           string     `gorm:"type:text;not null" json:"payload"`
	StatusCode        *int       `json:"status_code,omitempty"`
	AttemptCount      int        `gorm:"not null;default:0" json:"attempt_count"`
	DeliveredAt       *time.Time `json:"delivered_at,omitempty"`
	FailureReason     *string    `gorm:"type:text" json:"failure_reason,omitempty"`
}

type MessageLog struct {
	BaseModel
	TenantID          uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	WhatsAppSessionID uuid.UUID      `gorm:"type:uuid;index;not null" json:"whatsapp_session_id"`
	Direction         string         `gorm:"size:20;index;not null" json:"direction"`
	MessageType       string         `gorm:"size:30;not null" json:"message_type"`
	ToNumber          *string        `gorm:"size:40" json:"to_number,omitempty"`
	FromNumber        *string        `gorm:"size:40" json:"from_number,omitempty"`
	ExternalMessageID *string        `gorm:"size:120" json:"external_message_id,omitempty"`
	Status            *string        `gorm:"size:30" json:"status,omitempty"`
	Payload           datatypes.JSON `gorm:"type:jsonb" json:"payload"`
	ErrorMessage      *string        `gorm:"type:text" json:"error_message,omitempty"`
}

type AuditLog struct {
	BaseModel
	TenantID    *uuid.UUID     `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	ActorUserID *uuid.UUID     `gorm:"type:uuid;index" json:"actor_user_id,omitempty"`
	Action      string         `gorm:"size:80;index;not null" json:"action"`
	EntityType  string         `gorm:"size:60;not null" json:"entity_type"`
	EntityID    *uuid.UUID     `gorm:"type:uuid" json:"entity_id,omitempty"`
	RequestID   *string        `gorm:"size:80" json:"request_id,omitempty"`
	Metadata    datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
	OccurredAt  time.Time      `gorm:"not null" json:"occurred_at"`
}

func (m *BaseModel) BeforeCreate(_ *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
