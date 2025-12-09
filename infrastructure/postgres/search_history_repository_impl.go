package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type SearchHistoryRepositoryImpl struct {
	db *gorm.DB
}

func NewSearchHistoryRepository(db *gorm.DB) repositories.SearchHistoryRepository {
	return &SearchHistoryRepositoryImpl{db: db}
}

func (r *SearchHistoryRepositoryImpl) Create(ctx context.Context, history *models.SearchHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *SearchHistoryRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.SearchHistory, error) {
	var history models.SearchHistory
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

func (r *SearchHistoryRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.SearchHistory, error) {
	var histories []*models.SearchHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&histories).Error
	return histories, err
}

func (r *SearchHistoryRepositoryImpl) GetByUserIDAndType(ctx context.Context, userID uuid.UUID, searchType string, offset, limit int) ([]*models.SearchHistory, error) {
	var histories []*models.SearchHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND search_type = ?", userID, searchType).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&histories).Error
	return histories, err
}

func (r *SearchHistoryRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.SearchHistory{}).Error
}

func (r *SearchHistoryRepositoryImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.SearchHistory{}).Error
}

func (r *SearchHistoryRepositoryImpl) DeleteByUserIDAndType(ctx context.Context, userID uuid.UUID, searchType string) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND search_type = ?", userID, searchType).Delete(&models.SearchHistory{}).Error
}

func (r *SearchHistoryRepositoryImpl) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.SearchHistory{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *SearchHistoryRepositoryImpl) CountByUserIDAndType(ctx context.Context, userID uuid.UUID, searchType string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.SearchHistory{}).Where("user_id = ? AND search_type = ?", userID, searchType).Count(&count).Error
	return count, err
}

func (r *SearchHistoryRepositoryImpl) GetRecentByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*models.SearchHistory, error) {
	var histories []*models.SearchHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}
