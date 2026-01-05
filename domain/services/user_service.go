package services

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
)

type UserService interface {
	Register(ctx context.Context, req *dto.CreateUserRequest) (*models.User, error)
	Login(ctx context.Context, req *dto.LoginRequest) (string, *models.User, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	ListUsers(ctx context.Context, offset, limit int) ([]*models.User, int64, error)
	GenerateJWT(user *models.User) (string, error)
	ValidateJWT(token string) (*models.User, error)

	// Profile methods
	UpdateProfileInfo(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) (*models.User, error)
	UpdateAvatar(ctx context.Context, userID uuid.UUID, fileData []byte, contentType string) (*dto.UpdateAvatarResponse, error)
	DeleteAvatar(ctx context.Context, userID uuid.UUID) error

	// OAuth methods
	GetGoogleAuthURL(state string) string
	HandleGoogleCallback(ctx context.Context, code string, stateData dto.OAuthStateData) (string, *models.User, bool, error)
	GetLineAuthURL(state string) string
	HandleLineCallback(ctx context.Context, code string, stateData dto.OAuthStateData) (string, *models.User, bool, error)
}