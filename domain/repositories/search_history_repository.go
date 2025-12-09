package repositories

import (
	"context"

	"github.com/google/uuid"

	"gofiber-template/domain/models"
)

type SearchHistoryRepository interface {
	Create(ctx context.Context, history *models.SearchHistory) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.SearchHistory, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.SearchHistory, error)
	GetByUserIDAndType(ctx context.Context, userID uuid.UUID, searchType string, offset, limit int) ([]*models.SearchHistory, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteByUserIDAndType(ctx context.Context, userID uuid.UUID, searchType string) error
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	CountByUserIDAndType(ctx context.Context, userID uuid.UUID, searchType string) (int64, error)
	GetRecentByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*models.SearchHistory, error)
}
