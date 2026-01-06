package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type APIRequestLogRepositoryImpl struct {
	db *gorm.DB
}

func NewAPIRequestLogRepository(db *gorm.DB) repositories.APIRequestLogRepository {
	return &APIRequestLogRepositoryImpl{db: db}
}

func (r *APIRequestLogRepositoryImpl) Create(ctx context.Context, log *models.APIRequestLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *APIRequestLogRepositoryImpl) CreateBatch(ctx context.Context, logs []*models.APIRequestLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&logs).Error
}

func (r *APIRequestLogRepositoryImpl) GetStatsByService(ctx context.Context, startDate, endDate time.Time) ([]models.APIRequestStats, error) {
	var stats []models.APIRequestStats

	err := r.db.WithContext(ctx).
		Model(&models.APIRequestLog{}).
		Select(`
			service_name,
			'' as endpoint,
			COUNT(*) as total_requests,
			SUM(CASE WHEN source = 'cache' OR source = 'database' THEN 1 ELSE 0 END) as cache_hits,
			SUM(CASE WHEN source = 'api' THEN 1 ELSE 0 END) as api_hits,
			COALESCE(SUM(CASE WHEN source = 'cache' OR source = 'database' THEN 1 ELSE 0 END) * 100.0 / NULLIF(COUNT(*), 0), 0) as cache_hit_rate,
			COALESCE(SUM(estimated_cost), 0) as total_cost,
			COALESCE(AVG(duration_ms), 0) as avg_duration_ms
		`).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("service_name").
		Order("total_requests DESC").
		Scan(&stats).Error

	return stats, err
}

func (r *APIRequestLogRepositoryImpl) GetStatsByEndpoint(ctx context.Context, serviceName string, startDate, endDate time.Time) ([]models.APIRequestStats, error) {
	var stats []models.APIRequestStats

	query := r.db.WithContext(ctx).
		Model(&models.APIRequestLog{}).
		Select(`
			service_name,
			endpoint,
			COUNT(*) as total_requests,
			SUM(CASE WHEN source = 'cache' OR source = 'database' THEN 1 ELSE 0 END) as cache_hits,
			SUM(CASE WHEN source = 'api' THEN 1 ELSE 0 END) as api_hits,
			COALESCE(SUM(CASE WHEN source = 'cache' OR source = 'database' THEN 1 ELSE 0 END) * 100.0 / NULLIF(COUNT(*), 0), 0) as cache_hit_rate,
			COALESCE(SUM(estimated_cost), 0) as total_cost,
			COALESCE(AVG(duration_ms), 0) as avg_duration_ms
		`).
		Where("created_at BETWEEN ? AND ?", startDate, endDate)

	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	err := query.
		Group("service_name, endpoint").
		Order("total_requests DESC").
		Scan(&stats).Error

	return stats, err
}

func (r *APIRequestLogRepositoryImpl) GetDailyStats(ctx context.Context, startDate, endDate time.Time) ([]models.DailyStats, error) {
	var stats []models.DailyStats

	err := r.db.WithContext(ctx).
		Model(&models.APIRequestLog{}).
		Select(`
			TO_CHAR(created_at, 'YYYY-MM-DD') as date,
			COUNT(*) as total_requests,
			SUM(CASE WHEN source = 'cache' OR source = 'database' THEN 1 ELSE 0 END) as cache_hits,
			SUM(CASE WHEN source = 'api' THEN 1 ELSE 0 END) as api_hits,
			COALESCE(SUM(estimated_cost), 0) as total_cost
		`).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("TO_CHAR(created_at, 'YYYY-MM-DD')").
		Order("date DESC").
		Scan(&stats).Error

	return stats, err
}

func (r *APIRequestLogRepositoryImpl) GetServiceCosts(ctx context.Context, startDate, endDate time.Time) ([]models.ServiceCost, error) {
	var costs []models.ServiceCost

	err := r.db.WithContext(ctx).
		Model(&models.APIRequestLog{}).
		Select(`
			service_name,
			COUNT(*) as total_requests,
			COALESCE(SUM(estimated_cost), 0) as total_cost,
			COALESCE(SUM(estimated_cost) / NULLIF(COUNT(*), 0), 0) as cost_per_request
		`).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("source = ?", "api"). // Only count actual API calls
		Group("service_name").
		Order("total_cost DESC").
		Scan(&costs).Error

	return costs, err
}

func (r *APIRequestLogRepositoryImpl) GetTotalCost(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	var totalCost float64

	err := r.db.WithContext(ctx).
		Model(&models.APIRequestLog{}).
		Select("COALESCE(SUM(estimated_cost), 0)").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("source = ?", "api").
		Scan(&totalCost).Error

	return totalCost, err
}

func (r *APIRequestLogRepositoryImpl) GetCacheHitRate(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	var result struct {
		CacheHits int64
		Total     int64
	}

	err := r.db.WithContext(ctx).
		Model(&models.APIRequestLog{}).
		Select(`
			SUM(CASE WHEN source = 'cache' OR source = 'database' THEN 1 ELSE 0 END) as cache_hits,
			COUNT(*) as total
		`).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	if result.Total == 0 {
		return 0, nil
	}

	return float64(result.CacheHits) / float64(result.Total) * 100, nil
}

func (r *APIRequestLogRepositoryImpl) CountBySource(ctx context.Context, source string, startDate, endDate time.Time) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.APIRequestLog{}).
		Where("source = ?", source).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error

	return count, err
}

func (r *APIRequestLogRepositoryImpl) GetRecentLogs(ctx context.Context, limit int) ([]*models.APIRequestLog, error) {
	var logs []*models.APIRequestLog

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error

	return logs, err
}

func (r *APIRequestLogRepositoryImpl) DeleteOldLogs(ctx context.Context, before time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("created_at < ?", before).
		Delete(&models.APIRequestLog{})

	return result.RowsAffected, result.Error
}
