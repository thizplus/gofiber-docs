package services

import (
	"context"

	"github.com/google/uuid"

	"gofiber-template/domain/dto"
)

type SearchService interface {
	// Unified search
	Search(ctx context.Context, userID uuid.UUID, req *dto.SearchRequest) (*dto.SearchResponse, error)

	// Website search
	SearchWebsites(ctx context.Context, userID uuid.UUID, req *dto.SearchRequest) (*dto.WebsiteSearchResponse, error)

	// Image search
	SearchImages(ctx context.Context, userID uuid.UUID, req *dto.ImageSearchRequest) (*dto.ImageSearchResponse, error)

	// Video search (YouTube)
	SearchVideos(ctx context.Context, userID uuid.UUID, req *dto.VideoSearchRequest) (*dto.VideoSearchResponse, error)
	GetVideoDetails(ctx context.Context, videoID string) (*dto.VideoResult, error)

	// Place search (Google Places)
	SearchPlaces(ctx context.Context, userID uuid.UUID, req *dto.PlaceSearchRequest) (*dto.PlaceSearchResponse, error)
	GetPlaceDetails(ctx context.Context, placeID string, userLat, userLng float64) (*dto.PlaceDetailResponse, error)
	GetPlaceDetailsEnhanced(ctx context.Context, placeID string, userLat, userLng float64, includeAI bool) (*dto.PlaceDetailEnhancedResponse, error)
	SearchNearbyPlaces(ctx context.Context, req *dto.NearbyPlacesRequest) (*dto.PlaceSearchResponse, error)

	// Search history
	GetSearchHistory(ctx context.Context, userID uuid.UUID, req *dto.GetSearchHistoryRequest) (*dto.SearchHistoryListResponse, error)
	ClearSearchHistory(ctx context.Context, userID uuid.UUID, req *dto.ClearSearchHistoryRequest) error
	DeleteSearchHistoryItem(ctx context.Context, userID uuid.UUID, historyID uuid.UUID) error
}
