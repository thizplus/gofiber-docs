package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// AIChatSession represents an AI chat session
type AIChatSession struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Title        string    `gorm:"type:varchar(255)"`
	InitialQuery string    `gorm:"type:varchar(500)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Relationships
	User     User            `gorm:"foreignKey:UserID"`
	Messages []AIChatMessage `gorm:"foreignKey:SessionID"`
}

func (AIChatSession) TableName() string {
	return "ai_chat_sessions"
}

// AIChatMessage represents a message in an AI chat session
type AIChatMessage struct {
	ID        uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SessionID uuid.UUID      `gorm:"type:uuid;not null;index"`
	Role      string         `gorm:"type:varchar(20);not null"` // user, assistant
	Content   string         `gorm:"type:text;not null"`
	Sources   datatypes.JSON `gorm:"type:jsonb;default:'[]'"`
	CreatedAt time.Time

	// Relationships
	Session AIChatSession `gorm:"foreignKey:SessionID"`
}

func (AIChatMessage) TableName() string {
	return "ai_chat_messages"
}

// Message roles
const (
	MessageRoleUser      = "user"
	MessageRoleAssistant = "assistant"
)

// MessageSource represents a source in AI response
type MessageSource struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet,omitempty"`
}
