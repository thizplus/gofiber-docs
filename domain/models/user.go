package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Username  string    `gorm:"uniqueIndex;not null"`
	Password  string    // Empty for OAuth users
	FirstName string
	LastName  string
	Avatar    string `gorm:"type:varchar(500)"` // Avatar URL from OAuth provider or R2

	// STOU Specific
	StudentID string `gorm:"type:varchar(11);uniqueIndex"` // รหัสนักศึกษา 11 หลัก

	// Preferences
	Language string `gorm:"type:varchar(5);default:'th'"`    // th, en
	Theme    string `gorm:"type:varchar(10);default:'light'"` // light, dark

	// Status
	Role     string `gorm:"default:'user'"`
	IsActive bool   `gorm:"default:true"`

	// OAuth Fields
	GoogleID     *string `gorm:"type:varchar(255);uniqueIndex"` // Google OAuth ID
	LineID       *string `gorm:"type:varchar(255);uniqueIndex"` // LINE OAuth ID
	AuthProvider string  `gorm:"type:varchar(20);default:'local'"` // local, google, line

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "users"
}

// IsOAuthUser returns true if user registered via OAuth
func (u *User) IsOAuthUser() bool {
	return u.AuthProvider != "local" && u.AuthProvider != ""
}