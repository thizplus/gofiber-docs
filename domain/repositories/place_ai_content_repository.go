package repositories

import (
	"context"

	"github.com/google/uuid"

	"gofiber-template/domain/models"
)

type PlaceAIContentRepository interface {
	// Create creates a new place AI content record
	Create(ctx context.Context, content *models.PlaceAIContent) error

	// GetByPlaceID gets content by Google Place ID
	GetByPlaceID(ctx context.Context, placeID string) (*models.PlaceAIContent, error)

	// GetByPlaceIDAndLanguage gets content by place ID and language
	GetByPlaceIDAndLanguage(ctx context.Context, placeID, language string) (*models.PlaceAIContent, error)

	// Update updates an existing record
	Update(ctx context.Context, content *models.PlaceAIContent) error

	// Upsert creates or updates a record
	Upsert(ctx context.Context, content *models.PlaceAIContent) error

	// Delete deletes a record by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByPlaceID deletes by place ID
	DeleteByPlaceID(ctx context.Context, placeID string) error

	// DeleteExpired deletes all expired records
	DeleteExpired(ctx context.Context) (int64, error)

	// Exists checks if content exists for a place
	Exists(ctx context.Context, placeID string) (bool, error)
}
