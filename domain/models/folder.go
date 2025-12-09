package models

import (
	"time"

	"github.com/google/uuid"
)

// Folder represents a user's folder for saving items
type Folder struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Name          string    `gorm:"type:varchar(255);not null"`
	Description   string    `gorm:"type:text"`
	CoverImageURL string    `gorm:"type:text"`
	IsPublic      bool      `gorm:"default:false"`
	ItemCount     int       `gorm:"default:0"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// Relationships
	User  User         `gorm:"foreignKey:UserID"`
	Items []FolderItem `gorm:"foreignKey:FolderID"`
}

func (Folder) TableName() string {
	return "folders"
}

// Folder item types
const (
	FolderItemTypePlace   = "place"
	FolderItemTypeWebsite = "website"
	FolderItemTypeImage   = "image"
	FolderItemTypeVideo   = "video"
	FolderItemTypeLink    = "link"
)
