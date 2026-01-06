package models

import (
	"time"

	"github.com/google/uuid"
)

// APIRequestLog tracks external API calls for cost monitoring
type APIRequestLog struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`

	// Request info
	ServiceName string `gorm:"type:varchar(50);not null;index"` // google_places, google_translate, youtube, openai
	Endpoint    string `gorm:"type:varchar(100);not null"`      // text_search, nearby_search, place_details, etc.

	// Source tracking
	Source      string `gorm:"type:varchar(20);not null;index"` // api, cache, database
	CacheKey    string `gorm:"type:varchar(255)"`               // Cache key if from cache

	// Request details
	RequestParams string `gorm:"type:text"`     // JSON of request parameters (for debugging)
	ResponseSize  int    `gorm:"default:0"`     // Response size in bytes

	// Cost estimation
	EstimatedCost float64 `gorm:"type:decimal(10,6);default:0"` // Estimated cost in USD
	FieldsUsed    string  `gorm:"type:varchar(500)"`            // Fields requested (for Places API)

	// User tracking
	UserID    *uuid.UUID `gorm:"type:uuid;index"` // NULL for guest
	IPAddress string     `gorm:"type:varchar(45)"`
	UserAgent string     `gorm:"type:varchar(500)"`

	// Status
	Success      bool   `gorm:"default:true"`
	ErrorMessage string `gorm:"type:text"`

	// Performance
	DurationMs int `gorm:"default:0"` // Request duration in milliseconds

	// Timestamps
	CreatedAt time.Time `gorm:"index"`
}

func (APIRequestLog) TableName() string {
	return "api_request_logs"
}

// APIRequestStats represents aggregated statistics
type APIRequestStats struct {
	ServiceName   string  `json:"serviceName"`
	Endpoint      string  `json:"endpoint"`
	TotalRequests int64   `json:"totalRequests"`
	CacheHits     int64   `json:"cacheHits"`
	APIHits       int64   `json:"apiHits"`
	CacheHitRate  float64 `json:"cacheHitRate"`
	TotalCost     float64 `json:"totalCost"`
	AvgDurationMs float64 `json:"avgDurationMs"`
}

// DailyStats represents daily statistics
type DailyStats struct {
	Date          string  `json:"date"`
	TotalRequests int64   `json:"totalRequests"`
	CacheHits     int64   `json:"cacheHits"`
	APIHits       int64   `json:"apiHits"`
	TotalCost     float64 `json:"totalCost"`
}

// ServiceCost represents cost breakdown by service
type ServiceCost struct {
	ServiceName   string  `json:"serviceName"`
	TotalRequests int64   `json:"totalRequests"`
	TotalCost     float64 `json:"totalCost"`
	CostPerRequest float64 `json:"costPerRequest"`
}

// Estimated costs per API call (USD)
const (
	CostPlacesTextSearch   = 0.032  // $32 per 1000 requests
	CostPlacesNearbySearch = 0.032  // $32 per 1000 requests
	CostPlacesDetailsBasic = 0.0    // Free
	CostPlacesDetailsContact = 0.003 // $3 per 1000
	CostPlacesDetailsAtmosphere = 0.005 // $5 per 1000
	CostPlacesPhoto        = 0.007  // $7 per 1000
	CostGoogleTranslate    = 0.00002 // $20 per 1M characters
	CostYouTubeSearch      = 0.0001 // Quota based
	CostOpenAI             = 0.002  // ~$2 per 1000 tokens (varies)
)
