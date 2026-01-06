package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// PlaceAIContent stores AI-generated content for places
// This prevents repeated API calls for the same place
// Each place can have content in multiple languages
type PlaceAIContent struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PlaceID   string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_place_lang"` // Google Place ID
	PlaceName string    `gorm:"type:varchar(500);not null"`

	// AI Overview
	Summary        string         `gorm:"type:text"`
	History        string         `gorm:"type:text"`
	Highlights     datatypes.JSON `gorm:"type:jsonb;default:'[]'"` // []string
	BestTimeToVisit string        `gorm:"type:text"`
	Tips           datatypes.JSON `gorm:"type:jsonb;default:'[]'"` // []string

	// Guide Info
	QuickFacts      datatypes.JSON `gorm:"type:jsonb;default:'[]'"` // []string
	TalkingPoints   datatypes.JSON `gorm:"type:jsonb;default:'[]'"` // []string
	CommonQuestions datatypes.JSON `gorm:"type:jsonb;default:'[]'"` // []FAQ

	// Related Videos (cached YouTube results)
	RelatedVideos datatypes.JSON `gorm:"type:jsonb;default:'[]'"` // []VideoInfo

	// Metadata
	Language    string    `gorm:"type:varchar(10);default:'th';uniqueIndex:idx_place_lang"`
	GeneratedAt time.Time `gorm:"not null"`
	ExpiresAt   time.Time `gorm:"not null;index"` // For cleanup old records
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (PlaceAIContent) TableName() string {
	return "place_ai_contents"
}

// FAQ represents a frequently asked question
type FAQ struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// VideoInfo represents a related video
type VideoInfo struct {
	VideoID      string `json:"videoId"`
	Title        string `json:"title"`
	ThumbnailURL string `json:"thumbnailUrl"`
	ChannelTitle string `json:"channelTitle"`
	Duration     string `json:"duration,omitempty"`
	ViewCount    string `json:"viewCount,omitempty"`
}

// AIOverview represents the AI-generated overview
type AIOverview struct {
	Summary         string   `json:"summary"`
	History         string   `json:"history"`
	Highlights      []string `json:"highlights"`
	BestTimeToVisit string   `json:"bestTimeToVisit"`
	Tips            []string `json:"tips"`
	GeneratedAt     string   `json:"generatedAt"`
}

// GuideInfo represents info for tour guides
type GuideInfo struct {
	QuickFacts      []string `json:"quickFacts"`
	TalkingPoints   []string `json:"talkingPoints"`
	CommonQuestions []FAQ    `json:"commonQuestions"`
}
