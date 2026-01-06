package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type PlaceAIContentRepositoryImpl struct {
	db *gorm.DB
}

func NewPlaceAIContentRepository(db *gorm.DB) repositories.PlaceAIContentRepository {
	return &PlaceAIContentRepositoryImpl{db: db}
}

func (r *PlaceAIContentRepositoryImpl) Create(ctx context.Context, content *models.PlaceAIContent) error {
	return r.db.WithContext(ctx).Create(content).Error
}

func (r *PlaceAIContentRepositoryImpl) GetByPlaceID(ctx context.Context, placeID string) (*models.PlaceAIContent, error) {
	var content models.PlaceAIContent
	err := r.db.WithContext(ctx).
		Where("place_id = ?", placeID).
		Where("expires_at > ?", time.Now()).
		First(&content).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *PlaceAIContentRepositoryImpl) GetByPlaceIDAndLanguage(ctx context.Context, placeID, language string) (*models.PlaceAIContent, error) {
	var content models.PlaceAIContent
	err := r.db.WithContext(ctx).
		Where("place_id = ? AND language = ?", placeID, language).
		Where("expires_at > ?", time.Now()).
		First(&content).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *PlaceAIContentRepositoryImpl) Update(ctx context.Context, content *models.PlaceAIContent) error {
	return r.db.WithContext(ctx).Save(content).Error
}

func (r *PlaceAIContentRepositoryImpl) Upsert(ctx context.Context, content *models.PlaceAIContent) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "place_id"}, {Name: "language"}},
			UpdateAll: true,
		}).
		Create(content).Error
}

func (r *PlaceAIContentRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.PlaceAIContent{}).Error
}

func (r *PlaceAIContentRepositoryImpl) DeleteByPlaceID(ctx context.Context, placeID string) error {
	return r.db.WithContext(ctx).Where("place_id = ?", placeID).Delete(&models.PlaceAIContent{}).Error
}

func (r *PlaceAIContentRepositoryImpl) DeleteExpired(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&models.PlaceAIContent{})
	return result.RowsAffected, result.Error
}

func (r *PlaceAIContentRepositoryImpl) Exists(ctx context.Context, placeID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.PlaceAIContent{}).
		Where("place_id = ?", placeID).
		Where("expires_at > ?", time.Now()).
		Count(&count).Error
	return count > 0, err
}
