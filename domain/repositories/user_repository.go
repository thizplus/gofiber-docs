package repositories

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, id uuid.UUID, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.User, error)
	Count(ctx context.Context) (int64, error)

	// OAuth methods
	GetByGoogleID(ctx context.Context, googleID string) (*models.User, error)
	GetByLineID(ctx context.Context, lineID string) (*models.User, error)

	// Profile methods
	GetByStudentID(ctx context.Context, studentID string) (*models.User, error)
}