package serviceimpl

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"

	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

// APILoggerService handles logging of external API requests
type APILoggerService struct {
	repo    repositories.APIRequestLogRepository
	buffer  []*models.APIRequestLog
	mu      sync.Mutex
	maxSize int
}

// NewAPILoggerService creates a new API logger service
func NewAPILoggerService(repo repositories.APIRequestLogRepository) *APILoggerService {
	service := &APILoggerService{
		repo:    repo,
		buffer:  make([]*models.APIRequestLog, 0, 100),
		maxSize: 100, // Flush every 100 logs
	}

	// Start background flusher
	go service.startFlusher()

	return service
}

// LogRequest logs an API request
func (s *APILoggerService) LogRequest(ctx context.Context, log *models.APIRequestLog) {
	log.CreatedAt = time.Now()

	s.mu.Lock()
	s.buffer = append(s.buffer, log)
	shouldFlush := len(s.buffer) >= s.maxSize
	s.mu.Unlock()

	if shouldFlush {
		go s.flush(ctx)
	}
}

// LogAPICall logs an external API call (source = "api")
func (s *APILoggerService) LogAPICall(ctx context.Context, serviceName, endpoint string, params interface{}, cost float64, durationMs int, userID *uuid.UUID, success bool, errMsg string) {
	paramsJSON := ""
	if params != nil {
		if b, err := json.Marshal(params); err == nil {
			paramsJSON = string(b)
		}
	}

	log := &models.APIRequestLog{
		ServiceName:   serviceName,
		Endpoint:      endpoint,
		Source:        "api",
		RequestParams: paramsJSON,
		EstimatedCost: cost,
		DurationMs:    durationMs,
		UserID:        userID,
		Success:       success,
		ErrorMessage:  errMsg,
	}

	s.LogRequest(ctx, log)
}

// LogCacheHit logs a cache hit (source = "cache")
func (s *APILoggerService) LogCacheHit(ctx context.Context, serviceName, endpoint, cacheKey string, userID *uuid.UUID) {
	log := &models.APIRequestLog{
		ServiceName:   serviceName,
		Endpoint:      endpoint,
		Source:        "cache",
		CacheKey:      cacheKey,
		EstimatedCost: 0,
		UserID:        userID,
		Success:       true,
	}

	s.LogRequest(ctx, log)
}

// LogDatabaseHit logs a database cache hit (source = "database")
func (s *APILoggerService) LogDatabaseHit(ctx context.Context, serviceName, endpoint string, userID *uuid.UUID) {
	log := &models.APIRequestLog{
		ServiceName:   serviceName,
		Endpoint:      endpoint,
		Source:        "database",
		EstimatedCost: 0,
		UserID:        userID,
		Success:       true,
	}

	s.LogRequest(ctx, log)
}

// flush writes buffered logs to database
func (s *APILoggerService) flush(ctx context.Context) {
	s.mu.Lock()
	if len(s.buffer) == 0 {
		s.mu.Unlock()
		return
	}

	logs := s.buffer
	s.buffer = make([]*models.APIRequestLog, 0, 100)
	s.mu.Unlock()

	// Write to database
	if err := s.repo.CreateBatch(ctx, logs); err != nil {
		// Log error but don't fail - we don't want logging to break the app
		// In production, you might want to write to a fallback location
	}
}

// startFlusher starts a background goroutine that flushes logs periodically
func (s *APILoggerService) startFlusher() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.flush(context.Background())
	}
}

// Flush forces a flush of the buffer
func (s *APILoggerService) Flush(ctx context.Context) {
	s.flush(ctx)
}

// GetStatsByService returns statistics grouped by service
func (s *APILoggerService) GetStatsByService(ctx context.Context, days int) ([]models.APIRequestStats, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)
	return s.repo.GetStatsByService(ctx, startDate, endDate)
}

// GetStatsByEndpoint returns statistics grouped by endpoint
func (s *APILoggerService) GetStatsByEndpoint(ctx context.Context, serviceName string, days int) ([]models.APIRequestStats, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)
	return s.repo.GetStatsByEndpoint(ctx, serviceName, startDate, endDate)
}

// GetDailyStats returns daily statistics
func (s *APILoggerService) GetDailyStats(ctx context.Context, days int) ([]models.DailyStats, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)
	return s.repo.GetDailyStats(ctx, startDate, endDate)
}

// GetServiceCosts returns cost breakdown by service
func (s *APILoggerService) GetServiceCosts(ctx context.Context, days int) ([]models.ServiceCost, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)
	return s.repo.GetServiceCosts(ctx, startDate, endDate)
}

// GetTotalCost returns total estimated cost
func (s *APILoggerService) GetTotalCost(ctx context.Context, days int) (float64, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)
	return s.repo.GetTotalCost(ctx, startDate, endDate)
}

// GetCacheHitRate returns cache hit rate percentage
func (s *APILoggerService) GetCacheHitRate(ctx context.Context, days int) (float64, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)
	return s.repo.GetCacheHitRate(ctx, startDate, endDate)
}

// GetSummary returns a summary of API usage
func (s *APILoggerService) GetSummary(ctx context.Context, days int) (*APISummary, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	stats, err := s.repo.GetStatsByService(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	costs, err := s.repo.GetServiceCosts(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	cacheHitRate, err := s.repo.GetCacheHitRate(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	totalCost, err := s.repo.GetTotalCost(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var totalRequests, apiHits, cacheHits int64
	for _, stat := range stats {
		totalRequests += stat.TotalRequests
		apiHits += stat.APIHits
		cacheHits += stat.CacheHits
	}

	return &APISummary{
		Period:        days,
		TotalRequests: totalRequests,
		APIHits:       apiHits,
		CacheHits:     cacheHits,
		CacheHitRate:  cacheHitRate,
		TotalCost:     totalCost,
		CostSaved:     totalCost * (cacheHitRate / 100), // Estimated savings from cache
		ServiceStats:  stats,
		ServiceCosts:  costs,
	}, nil
}

// CleanupOldLogs removes logs older than specified days
func (s *APILoggerService) CleanupOldLogs(ctx context.Context, days int) (int64, error) {
	before := time.Now().AddDate(0, 0, -days)
	return s.repo.DeleteOldLogs(ctx, before)
}

// APISummary represents a summary of API usage
type APISummary struct {
	Period        int                    `json:"period"` // Days
	TotalRequests int64                  `json:"totalRequests"`
	APIHits       int64                  `json:"apiHits"`
	CacheHits     int64                  `json:"cacheHits"`
	CacheHitRate  float64                `json:"cacheHitRate"`
	TotalCost     float64                `json:"totalCost"`
	CostSaved     float64                `json:"costSaved"`
	ServiceStats  []models.APIRequestStats `json:"serviceStats"`
	ServiceCosts  []models.ServiceCost     `json:"serviceCosts"`
}
