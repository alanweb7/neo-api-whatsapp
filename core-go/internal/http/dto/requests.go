package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	TenantID string `json:"tenant_id" binding:"required,uuid"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	TenantID     string `json:"tenant_id" binding:"required,uuid"`
}

type CreateTenantRequest struct {
	Name string `json:"name" binding:"required,min=3,max=120"`
	Slug string `json:"slug" binding:"required,min=3,max=80"`
}

type UpdateTenantRequest struct {
	Name   *string `json:"name"`
	Status *string `json:"status"`
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required,min=2,max=120"`
	Role     string `json:"role"`
}

type AttachUserRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
	Role   string `json:"role"`
}

type CreateAPIKeyRequest struct {
	Name string `json:"name" binding:"required,min=3,max=80"`
}

type CreateSessionRequest struct {
	Name string `json:"name" binding:"required,min=2,max=120"`
}

type SendTextRequest struct {
	SessionID string `json:"session_id" binding:"required,uuid"`
	To        string `json:"to" binding:"required"`
	Text      string `json:"text" binding:"required"`
}

type SendMediaRequest struct {
	SessionID string `json:"session_id" binding:"required,uuid"`
	To        string `json:"to" binding:"required"`
	MediaURL  string `json:"media_url" binding:"required,url"`
	Caption   string `json:"caption"`
	FileName  string `json:"file_name"`
}

type ButtonItem struct {
	Type        string `json:"type" binding:"required,oneof=quick_reply"`
	DisplayText string `json:"displayText" binding:"required,min=1,max=40"`
	ID          string `json:"id" binding:"required,min=1,max=128"`
}

type SendButtonsRequest struct {
	SessionID    string       `json:"session_id" binding:"required,uuid"`
	JID          string       `json:"jid"`
	To           string       `json:"to"`
	Text         string       `json:"text" binding:"required,min=1,max=1024"`
	Footer       string       `json:"footer"`
	Buttons      []ButtonItem `json:"buttons" binding:"required,min=1,max=3,dive"`
	FallbackText string       `json:"fallback_text"`
}

type CreateWebhookRequest struct {
	Name       string   `json:"name" binding:"required,min=3,max=80"`
	URL        string   `json:"url" binding:"required,url"`
	Secret     string   `json:"secret"`
	EventTypes []string `json:"event_types" binding:"required,min=1"`
}

type UpdateWebhookRequest struct {
	Name       *string  `json:"name"`
	URL        *string  `json:"url"`
	Secret     *string  `json:"secret"`
	IsActive   *bool    `json:"is_active"`
	EventTypes []string `json:"event_types"`
}
