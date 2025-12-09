package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// FolderItem represents an item saved in a folder
type FolderItem struct {
	ID           uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FolderID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	Type         string         `gorm:"type:varchar(50);not null"` // place, website, image, video, link
	Title        string         `gorm:"type:varchar(255);not null"`
	URL          string         `gorm:"type:text;not null"`
	ThumbnailURL string         `gorm:"type:text"`
	Description  string         `gorm:"type:text"`
	Metadata     datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	SortOrder    int            `gorm:"default:0"`
	CreatedAt    time.Time

	// Relationships
	Folder Folder `gorm:"foreignKey:FolderID"`
}

func (FolderItem) TableName() string {
	return "folder_items"
}

// PlaceMetadata contains place-specific metadata
type PlaceMetadata struct {
	PlaceID     string  `json:"place_id,omitempty"`
	Address     string  `json:"address,omitempty"`
	Rating      float64 `json:"rating,omitempty"`
	ReviewCount int     `json:"review_count,omitempty"`
	Lat         float64 `json:"lat,omitempty"`
	Lng         float64 `json:"lng,omitempty"`
}

// WebsiteMetadata contains website-specific metadata
type WebsiteMetadata struct {
	Snippet string `json:"snippet,omitempty"`
	Source  string `json:"source,omitempty"`
}

// ImageMetadata contains image-specific metadata
type ImageMetadata struct {
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
	Source string `json:"source,omitempty"`
}

// VideoMetadata contains video-specific metadata
type VideoMetadata struct {
	Duration  string `json:"duration,omitempty"`
	Channel   string `json:"channel,omitempty"`
	VideoID   string `json:"video_id,omitempty"`
	ViewCount int64  `json:"view_count,omitempty"`
}
