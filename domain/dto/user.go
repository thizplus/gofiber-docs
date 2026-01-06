package dto

import (
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email,max=255"`
	Username  string `json:"username" validate:"omitempty,min=3,max=20,alphanum"` // Optional: auto-generated from email if not provided
	Password  string `json:"password" validate:"required,min=8,max=72"`
	FirstName string `json:"firstName" validate:"required,min=1,max=50"`
	LastName  string `json:"lastName" validate:"required,min=1,max=50"`
}

type UpdateUserRequest struct {
	FirstName string `json:"firstName" validate:"omitempty,min=1,max=50"`
	LastName  string `json:"lastName" validate:"omitempty,min=1,max=50"`
	Avatar    string `json:"avatar" validate:"omitempty,url,max=500"`
}

// UpdateProfileRequest - request to update user profile
type UpdateProfileRequest struct {
	FirstName string `json:"firstName" validate:"omitempty,min=1,max=50"`
	LastName  string `json:"lastName" validate:"omitempty,min=1,max=50"`
	StudentID string `json:"studentId" validate:"omitempty,len=11,numeric"` // รหัสนักศึกษา 11 หลัก
	Language  string `json:"language" validate:"omitempty,oneof=th en"`
	Theme     string `json:"theme" validate:"omitempty,oneof=light dark"`
}

// UpdateAvatarResponse - response after avatar upload
type UpdateAvatarResponse struct {
	AvatarURL string `json:"avatarUrl"`
}

type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Avatar       string    `json:"avatar"`
	StudentID    string    `json:"studentId,omitempty"`
	Language     string    `json:"language"`
	Theme        string    `json:"theme"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"isActive"`
	AuthProvider string    `json:"authProvider"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Meta  PaginationMeta `json:"meta"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=NewPassword"`
}

// UserToUserResponse converts User model to UserResponse DTO
func UserToUserResponse(user *models.User) *UserResponse {
	if user == nil {
		return nil
	}

	// Handle nullable StudentID
	studentID := ""
	if user.StudentID != nil {
		studentID = *user.StudentID
	}

	return &UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Avatar:       user.Avatar,
		StudentID:    studentID,
		Language:     user.Language,
		Theme:        user.Theme,
		Role:         user.Role,
		IsActive:     user.IsActive,
		AuthProvider: user.AuthProvider,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
