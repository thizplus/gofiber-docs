package repositories

import (
	"context"

	"github.com/google/uuid"

	"gofiber-template/domain/models"
)

type FavoriteRepository interface {
	Create(ctx context.Context, favorite *models.Favorite) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Favorite, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Favorite, error)
	GetByUserIDAndType(ctx context.Context, userID uuid.UUID, favType string, offset, limit int) ([]*models.Favorite, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserIDAndURL(ctx context.Context, userID uuid.UUID, url string) error
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	CountByUserIDAndType(ctx context.Context, userID uuid.UUID, favType string) (int64, error)
	ExistsByUserIDAndURL(ctx context.Context, userID uuid.UUID, url string) (bool, error)
	ExistsByUserIDAndExternalID(ctx context.Context, userID uuid.UUID, externalID string) (bool, error)
	GetByUserIDAndURL(ctx context.Context, userID uuid.UUID, url string) (*models.Favorite, error)
	GetByUserIDAndExternalID(ctx context.Context, userID uuid.UUID, externalID string) (*models.Favorite, error)
}
