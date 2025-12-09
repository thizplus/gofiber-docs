package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Favorite represents a user's favorite item
type Favorite struct {
	ID           uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index"`
	Type         string         `gorm:"type:varchar(50);not null"` // place, website, image, video
	ExternalID   string         `gorm:"type:varchar(255);index"`   // Google Place ID, etc.
	Title        string         `gorm:"type:varchar(255);not null"`
	URL          string         `gorm:"type:text;not null"`
	ThumbnailURL string         `gorm:"type:text"`
	Rating       float64        `gorm:"type:decimal(2,1)"`
	ReviewCount  int            `gorm:"default:0"`
	Address      string         `gorm:"type:text"`
	Metadata     datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	CreatedAt    time.Time

	// Relationships
	User User `gorm:"foreignKey:UserID"`
}

func (Favorite) TableName() string {
	return "favorites"
}

// Favorite types
const (
	FavoriteTypePlace   = "place"
	FavoriteTypeWebsite = "website"
	FavoriteTypeImage   = "image"
	FavoriteTypeVideo   = "video"
)
