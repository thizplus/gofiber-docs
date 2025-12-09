package services

import (
	"context"

	"github.com/google/uuid"

	"gofiber-template/domain/dto"
)

type FavoriteService interface {
	// Favorite operations
	AddFavorite(ctx context.Context, userID uuid.UUID, req *dto.AddFavoriteRequest) (*dto.FavoriteResponse, error)
	GetFavorites(ctx context.Context, userID uuid.UUID, req *dto.GetFavoritesRequest) (*dto.FavoriteListResponse, error)
	RemoveFavorite(ctx context.Context, userID uuid.UUID, favoriteID uuid.UUID) error

	// Check favorite status
	CheckFavorite(ctx context.Context, userID uuid.UUID, req *dto.CheckFavoriteRequest) (*dto.CheckFavoriteResponse, error)

	// Toggle favorite (add if not exists, remove if exists)
	ToggleFavorite(ctx context.Context, userID uuid.UUID, req *dto.AddFavoriteRequest) (*dto.CheckFavoriteResponse, error)
}
