package models

import (
	"time"

	"github.com/google/uuid"
)

// SearchHistory represents a user's search history
type SearchHistory struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Query       string    `gorm:"type:varchar(500);not null"`
	SearchType  string    `gorm:"type:varchar(50);not null;default:'all'"` // all, website, image, video, map, ai
	ResultCount int       `gorm:"default:0"`
	CreatedAt   time.Time `gorm:"index"`

	// Relationships
	User User `gorm:"foreignKey:UserID"`
}

func (SearchHistory) TableName() string {
	return "search_history"
}

// Search types
const (
	SearchTypeAll     = "all"
	SearchTypeWebsite = "website"
	SearchTypeImage   = "image"
	SearchTypeVideo   = "video"
	SearchTypeMap     = "map"
	SearchTypeAI      = "ai"
)
