package dto

import (
	"time"

	"github.com/google/uuid"
)

// ==================== Search Request DTOs ====================

type SearchRequest struct {
	Query    string `json:"query" query:"q" validate:"required,min=1,max=500"`
	Type     string `json:"type" query:"type" validate:"omitempty,oneof=all website image video map ai"`
	Page     int    `json:"page" query:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"pageSize" query:"pageSize" validate:"omitempty,min=1,max=50"`
	Language string `json:"language" query:"lang" validate:"omitempty,len=2"`
}

type PlaceSearchRequest struct {
	Query     string  `json:"query" query:"q" validate:"required,min=1,max=500"`
	Lat       float64 `json:"lat" query:"lat" validate:"omitempty,latitude"`
	Lng       float64 `json:"lng" query:"lng" validate:"omitempty,longitude"`
	Radius    int     `json:"radius" query:"radius" validate:"omitempty,min=100,max=50000"`
	PlaceType string  `json:"placeType" query:"type" validate:"omitempty"`
	Page      int     `json:"page" query:"page" validate:"omitempty,min=1"`
	PageSize  int     `json:"pageSize" query:"pageSize" validate:"omitempty,min=1,max=20"`
}

type VideoSearchRequest struct {
	Query    string `json:"query" query:"q" validate:"required,min=1,max=500"`
	Page     int    `json:"page" query:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"pageSize" query:"pageSize" validate:"omitempty,min=1,max=50"`
	Order    string `json:"order" query:"order" validate:"omitempty,oneof=relevance date viewCount rating"`
}

type ImageSearchRequest struct {
	Query    string `json:"query" query:"q" validate:"required,min=1,max=500"`
	Page     int    `json:"page" query:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"pageSize" query:"pageSize" validate:"omitempty,min=1,max=50"`
	Size     string `json:"size" query:"size" validate:"omitempty,oneof=small medium large"`
}

// ==================== Search Response DTOs ====================

type SearchResponse struct {
	Query      string         `json:"query"`
	Type       string         `json:"type"`
	Results    []SearchResult `json:"results"`
	TotalCount int64          `json:"totalCount"`
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
}

type SearchResult struct {
	Type         string   `json:"type"` // website, image, video, place
	Title        string   `json:"title"`
	URL          string   `json:"url,omitempty"`
	Snippet      string   `json:"snippet,omitempty"`
	ThumbnailURL string   `json:"thumbnailUrl,omitempty"`
	Source       string   `json:"source,omitempty"`
	PublishedAt  string   `json:"publishedAt,omitempty"`
	Rating       float64  `json:"rating,omitempty"`
	ReviewCount  int      `json:"reviewCount,omitempty"`
	// Place fields
	PlaceID string   `json:"placeId,omitempty"`
	Lat     float64  `json:"lat,omitempty"`
	Lng     float64  `json:"lng,omitempty"`
	Types   []string `json:"types,omitempty"`
	// Video fields
	VideoID   string `json:"videoId,omitempty"`
	Duration  string `json:"duration,omitempty"`
	ViewCount int64  `json:"viewCount,omitempty"`
}

type WebsiteSearchResponse struct {
	Query      string          `json:"query"`
	Results    []WebsiteResult `json:"results"`
	TotalCount int64           `json:"totalCount"`
	Page       int             `json:"page"`
	PageSize   int             `json:"pageSize"`
}

type WebsiteResult struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Snippet     string `json:"snippet"`
	DisplayLink string `json:"displayLink"`
	FormattedAt string `json:"formattedAt,omitempty"`
}

type ImageSearchResponse struct {
	Query      string        `json:"query"`
	Results    []ImageResult `json:"results"`
	TotalCount int64         `json:"totalCount"`
	Page       int           `json:"page"`
	PageSize   int           `json:"pageSize"`
}

type ImageResult struct {
	Title        string `json:"title"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Source       string `json:"source"`
	ContextLink  string `json:"contextLink"`
}

type VideoSearchResponse struct {
	Query      string        `json:"query"`
	Results    []VideoResult `json:"results"`
	TotalCount int64         `json:"totalCount"`
	Page       int           `json:"page"`
	PageSize   int           `json:"pageSize"`
}

type VideoResult struct {
	VideoID      string `json:"videoId"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ThumbnailURL string `json:"thumbnailUrl"`
	ChannelTitle string `json:"channelTitle"`
	PublishedAt  string `json:"publishedAt"`
	Duration     string `json:"duration,omitempty"`
	ViewCount    int64  `json:"viewCount,omitempty"`
	LikeCount    int64  `json:"likeCount,omitempty"`
}

type PlaceSearchResponse struct {
	Query      string        `json:"query"`
	Results    []PlaceResult `json:"results"`
	TotalCount int64         `json:"totalCount"`
	Page       int           `json:"page"`
	PageSize   int           `json:"pageSize"`
}

type PlaceResult struct {
	PlaceID      string   `json:"placeId"`
	Name         string   `json:"name"`
	Address      string   `json:"address"`
	Lat          float64  `json:"lat"`
	Lng          float64  `json:"lng"`
	Rating       float64  `json:"rating"`
	ReviewCount  int      `json:"reviewCount"`
	PriceLevel   int      `json:"priceLevel,omitempty"`
	Types        []string `json:"types"`
	PhotoURL     string   `json:"photoUrl,omitempty"`
	IsOpen       *bool    `json:"isOpen,omitempty"`
	Distance     float64  `json:"distance,omitempty"` // distance in meters from user location
	DistanceText string   `json:"distanceText,omitempty"`
}

type PlaceDetailResponse struct {
	PlaceID          string          `json:"placeId"`
	Name             string          `json:"name"`
	FormattedAddress string          `json:"formattedAddress"`
	Lat              float64         `json:"lat"`
	Lng              float64         `json:"lng"`
	Rating           float64         `json:"rating"`
	ReviewCount      int             `json:"reviewCount"`
	PriceLevel       int             `json:"priceLevel,omitempty"`
	Types            []string        `json:"types"`
	Phone            string          `json:"phone,omitempty"`
	Website          string          `json:"website,omitempty"`
	GoogleMapsURL    string          `json:"googleMapsUrl"`
	OpeningHours     []string        `json:"openingHours,omitempty"`
	Reviews          []PlaceReview   `json:"reviews,omitempty"`
	Photos           []PlacePhoto    `json:"photos,omitempty"`
	Distance         float64         `json:"distance,omitempty"`
	DistanceText     string          `json:"distanceText,omitempty"`
}

type PlaceReview struct {
	Author    string `json:"author"`
	Rating    int    `json:"rating"`
	Text      string `json:"text"`
	Time      string `json:"time"`
	PhotoURL  string `json:"photoUrl,omitempty"`
}

type PlacePhoto struct {
	URL       string `json:"url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

// ==================== Search History DTOs ====================

type SearchHistoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Query       string    `json:"query"`
	SearchType  string    `json:"searchType"`
	ResultCount int       `json:"resultCount"`
	CreatedAt   time.Time `json:"createdAt"`
}

type SearchHistoryListResponse struct {
	Histories []SearchHistoryResponse `json:"histories"`
	Meta      PaginationMeta          `json:"meta"`
}

type GetSearchHistoryRequest struct {
	SearchType string `json:"searchType" query:"type" validate:"omitempty,oneof=all website image video map ai"`
	Page       int    `json:"page" query:"page" validate:"omitempty,min=1"`
	PageSize   int    `json:"pageSize" query:"pageSize" validate:"omitempty,min=1,max=50"`
}

type ClearSearchHistoryRequest struct {
	SearchType string `json:"searchType" validate:"omitempty,oneof=all website image video map ai"`
}
