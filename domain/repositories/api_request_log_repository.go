package repositories

import (
	"context"
	"time"

	"gofiber-template/domain/models"
)

type APIRequestLogRepository interface {
	// Create creates a new log entry
	Create(ctx context.Context, log *models.APIRequestLog) error

	// CreateBatch creates multiple log entries
	CreateBatch(ctx context.Context, logs []*models.APIRequestLog) error

	// GetStatsByService gets aggregated stats by service
	GetStatsByService(ctx context.Context, startDate, endDate time.Time) ([]models.APIRequestStats, error)

	// GetStatsByEndpoint gets aggregated stats by endpoint
	GetStatsByEndpoint(ctx context.Context, serviceName string, startDate, endDate time.Time) ([]models.APIRequestStats, error)

	// GetDailyStats gets daily statistics
	GetDailyStats(ctx context.Context, startDate, endDate time.Time) ([]models.DailyStats, error)

	// GetServiceCosts gets cost breakdown by service
	GetServiceCosts(ctx context.Context, startDate, endDate time.Time) ([]models.ServiceCost, error)

	// GetTotalCost gets total estimated cost for a period
	GetTotalCost(ctx context.Context, startDate, endDate time.Time) (float64, error)

	// GetCacheHitRate gets cache hit rate for a period
	GetCacheHitRate(ctx context.Context, startDate, endDate time.Time) (float64, error)

	// CountBySource counts requests by source (api, cache, database)
	CountBySource(ctx context.Context, source string, startDate, endDate time.Time) (int64, error)

	// GetRecentLogs gets recent log entries
	GetRecentLogs(ctx context.Context, limit int) ([]*models.APIRequestLog, error)

	// DeleteOldLogs deletes logs older than specified date
	DeleteOldLogs(ctx context.Context, before time.Time) (int64, error)
}
