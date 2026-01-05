package dto

import (
	"time"

	"github.com/google/uuid"
)

// ==================== Favorite DTOs ====================

type AddFavoriteRequest struct {
	Type         string                 `json:"type" validate:"required,oneof=place website image video"`
	ExternalID   string                 `json:"externalId" validate:"omitempty,max=255"` // Google Place ID, etc.
	Title        string                 `json:"title" validate:"required,min=1,max=255"`
	URL          string                 `json:"url" validate:"required,url,max=2000"`
	ThumbnailURL string                 `json:"thumbnailUrl" validate:"omitempty,url,max=2000"`
	Rating       float64                `json:"rating" validate:"omitempty,min=0,max=5"`
	ReviewCount  int                    `json:"reviewCount" validate:"omitempty,min=0"`
	Address      string                 `json:"address" validate:"omitempty,max=500"`
	Metadata     map[string]interface{} `json:"metadata" validate:"omitempty"`
}

type FavoriteResponse struct {
	ID           uuid.UUID              `json:"id"`
	Type         string                 `json:"type"`
	ExternalID   string                 `json:"externalId,omitempty"`
	Title        string                 `json:"title"`
	URL          string                 `json:"url"`
	ThumbnailURL string                 `json:"thumbnailUrl,omitempty"`
	Rating       float64                `json:"rating,omitempty"`
	ReviewCount  int                    `json:"reviewCount,omitempty"`
	Address      string                 `json:"address,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time              `json:"createdAt"`
}

type FavoriteListResponse struct {
	Favorites []FavoriteResponse `json:"favorites"`
	Meta      PaginationMeta     `json:"meta"`
}

type GetFavoritesRequest struct {
	Type     string `json:"type" query:"type" validate:"omitempty,oneof=place website image video"`
	Page     int    `json:"page" query:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"pageSize" query:"pageSize" validate:"omitempty,min=1,max=50"`
}

type CheckFavoriteRequest struct {
	Type       string `json:"type" query:"type" validate:"required,oneof=place website image video"`
	URL        string `json:"url" query:"url" validate:"required_without=ExternalID,omitempty,url"`
	ExternalID string `json:"externalId" query:"externalId" validate:"required_without=URL,omitempty"`
}

type CheckFavoriteResponse struct {
	IsFavorite bool       `json:"isFavorite"`
	FavoriteID *uuid.UUID `json:"favoriteId,omitempty"`
}

// ==================== Batch Check Favorite DTOs ====================

type BatchCheckFavoritesRequest struct {
	ExternalIDs []string `json:"externalIds" validate:"required,min=1,max=50"`
}

type BatchCheckFavoritesResponse struct {
	Items map[string]CheckFavoriteResponse `json:"items"`
}
