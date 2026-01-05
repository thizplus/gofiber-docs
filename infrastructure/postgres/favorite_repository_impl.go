package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type FavoriteRepositoryImpl struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) repositories.FavoriteRepository {
	return &FavoriteRepositoryImpl{db: db}
}

func (r *FavoriteRepositoryImpl) Create(ctx context.Context, favorite *models.Favorite) error {
	return r.db.WithContext(ctx).Create(favorite).Error
}

func (r *FavoriteRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Favorite, error) {
	var favorite models.Favorite
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&favorite).Error
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *FavoriteRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Favorite, error) {
	var favorites []*models.Favorite
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&favorites).Error
	return favorites, err
}

func (r *FavoriteRepositoryImpl) GetByUserIDAndType(ctx context.Context, userID uuid.UUID, favType string, offset, limit int) ([]*models.Favorite, error) {
	var favorites []*models.Favorite
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, favType).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&favorites).Error
	return favorites, err
}

func (r *FavoriteRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Favorite{}).Error
}

func (r *FavoriteRepositoryImpl) DeleteByUserIDAndURL(ctx context.Context, userID uuid.UUID, url string) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND url = ?", userID, url).Delete(&models.Favorite{}).Error
}

func (r *FavoriteRepositoryImpl) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Favorite{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *FavoriteRepositoryImpl) CountByUserIDAndType(ctx context.Context, userID uuid.UUID, favType string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Favorite{}).Where("user_id = ? AND type = ?", userID, favType).Count(&count).Error
	return count, err
}

func (r *FavoriteRepositoryImpl) ExistsByUserIDAndURL(ctx context.Context, userID uuid.UUID, url string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Favorite{}).
		Where("user_id = ? AND url = ?", userID, url).
		Count(&count).Error
	return count > 0, err
}

func (r *FavoriteRepositoryImpl) ExistsByUserIDAndExternalID(ctx context.Context, userID uuid.UUID, externalID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Favorite{}).
		Where("user_id = ? AND external_id = ?", userID, externalID).
		Count(&count).Error
	return count > 0, err
}

func (r *FavoriteRepositoryImpl) GetByUserIDAndURL(ctx context.Context, userID uuid.UUID, url string) (*models.Favorite, error) {
	var favorite models.Favorite
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND url = ?", userID, url).
		First(&favorite).Error
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *FavoriteRepositoryImpl) GetByUserIDAndExternalID(ctx context.Context, userID uuid.UUID, externalID string) (*models.Favorite, error) {
	var favorite models.Favorite
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND external_id = ?", userID, externalID).
		First(&favorite).Error
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *FavoriteRepositoryImpl) GetByUserIDAndExternalIDs(ctx context.Context, userID uuid.UUID, externalIDs []string) ([]*models.Favorite, error) {
	var favorites []*models.Favorite
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND external_id IN ?", userID, externalIDs).
		Find(&favorites).Error
	return favorites, err
}
