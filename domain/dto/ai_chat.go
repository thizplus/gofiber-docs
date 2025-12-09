package dto

import (
	"time"

	"github.com/google/uuid"
)

// ==================== AI Chat Session DTOs ====================

type CreateAIChatRequest struct {
	Query string `json:"query" validate:"required,min=1,max=500"`
}

type AIChatSessionResponse struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	InitialQuery string    `json:"initialQuery"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type AIChatSessionDetailResponse struct {
	ID           uuid.UUID             `json:"id"`
	Title        string                `json:"title"`
	InitialQuery string                `json:"initialQuery"`
	Messages     []AIChatMessageResponse `json:"messages"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
}

type AIChatSessionListResponse struct {
	Sessions []AIChatSessionResponse `json:"sessions"`
	Meta     PaginationMeta          `json:"meta"`
}

type GetAIChatSessionsRequest struct {
	Page     int `json:"page" query:"page" validate:"omitempty,min=1"`
	PageSize int `json:"pageSize" query:"pageSize" validate:"omitempty,min=1,max=50"`
}

// ==================== AI Chat Message DTOs ====================

type SendAIChatMessageRequest struct {
	SessionID uuid.UUID `json:"sessionId" param:"sessionId" validate:"required"`
	Message   string    `json:"message" validate:"required,min=1,max=2000"`
}

type AIChatMessageResponse struct {
	ID        uuid.UUID       `json:"id"`
	SessionID uuid.UUID       `json:"sessionId"`
	Role      string          `json:"role"` // user, assistant
	Content   string          `json:"content"`
	Sources   []MessageSource `json:"sources,omitempty"`
	CreatedAt time.Time       `json:"createdAt"`
}

type MessageSource struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet,omitempty"`
}

// ==================== AI Search (Quick AI) DTOs ====================

type AISearchRequest struct {
	Query    string `json:"query" query:"q" validate:"required,min=1,max=500"`
	Language string `json:"language" query:"lang" validate:"omitempty,len=2"`
}

type AISearchResponse struct {
	Query    string          `json:"query"`
	Summary  string          `json:"summary"`
	Sources  []MessageSource `json:"sources"`
	Keywords []string        `json:"keywords,omitempty"`
}

// ==================== Streaming Response DTOs ====================

type AIStreamChunk struct {
	Type    string `json:"type"` // content, source, done, error
	Content string `json:"content,omitempty"`
	Source  *MessageSource `json:"source,omitempty"`
	Error   string `json:"error,omitempty"`
}
