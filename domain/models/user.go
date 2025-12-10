package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StudentID *string   `gorm:"type:varchar(20);uniqueIndex"` // รหัสนักศึกษา มสธ. (NULL allowed for non-students)
	Email     string    `gorm:"uniqueIndex;not null"`
	Username  string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	FirstName string
	LastName  string
	Avatar    string
	Role      string    `gorm:"default:'user'"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "users"
}