package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type FolderRepositoryImpl struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) repositories.FolderRepository {
	return &FolderRepositoryImpl{db: db}
}

func (r *FolderRepositoryImpl) Create(ctx context.Context, folder *models.Folder) error {
	return r.db.WithContext(ctx).Create(folder).Error
}

func (r *FolderRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Folder, error) {
	var folder models.Folder
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&folder).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *FolderRepositoryImpl) GetByIDWithItems(ctx context.Context, id uuid.UUID) (*models.Folder, error) {
	var folder models.Folder
	err := r.db.WithContext(ctx).
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, created_at DESC")
		}).
		Where("id = ?", id).
		First(&folder).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *FolderRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Folder, error) {
	var folders []*models.Folder
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&folders).Error
	return folders, err
}

func (r *FolderRepositoryImpl) GetPublicFolders(ctx context.Context, offset, limit int) ([]*models.Folder, error) {
	var folders []*models.Folder
	err := r.db.WithContext(ctx).
		Where("is_public = ?", true).
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&folders).Error
	return folders, err
}

func (r *FolderRepositoryImpl) Update(ctx context.Context, id uuid.UUID, folder *models.Folder) error {
	// Use Select("*") to force update all fields including zero values like is_public=false
	return r.db.WithContext(ctx).Model(folder).Where("id = ?", id).Select("*").Updates(folder).Error
}

func (r *FolderRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete folder items first
		if err := tx.Where("folder_id = ?", id).Delete(&models.FolderItem{}).Error; err != nil {
			return err
		}
		// Delete folder
		return tx.Where("id = ?", id).Delete(&models.Folder{}).Error
	})
}

func (r *FolderRepositoryImpl) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Folder{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *FolderRepositoryImpl) IncrementItemCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Folder{}).
		Where("id = ?", id).
		UpdateColumn("item_count", gorm.Expr("item_count + ?", 1)).Error
}

func (r *FolderRepositoryImpl) DecrementItemCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Folder{}).
		Where("id = ? AND item_count > 0", id).
		UpdateColumn("item_count", gorm.Expr("item_count - ?", 1)).Error
}

// FolderItemRepositoryImpl
type FolderItemRepositoryImpl struct {
	db *gorm.DB
}

func NewFolderItemRepository(db *gorm.DB) repositories.FolderItemRepository {
	return &FolderItemRepositoryImpl{db: db}
}

func (r *FolderItemRepositoryImpl) Create(ctx context.Context, item *models.FolderItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *FolderItemRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.FolderItem, error) {
	var item models.FolderItem
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *FolderItemRepositoryImpl) GetByFolderID(ctx context.Context, folderID uuid.UUID, offset, limit int) ([]*models.FolderItem, error) {
	var items []*models.FolderItem
	err := r.db.WithContext(ctx).
		Where("folder_id = ?", folderID).
		Order("sort_order ASC, created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&items).Error
	return items, err
}

func (r *FolderItemRepositoryImpl) GetByFolderIDAndType(ctx context.Context, folderID uuid.UUID, itemType string, offset, limit int) ([]*models.FolderItem, error) {
	var items []*models.FolderItem
	err := r.db.WithContext(ctx).
		Where("folder_id = ? AND type = ?", folderID, itemType).
		Order("sort_order ASC, created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&items).Error
	return items, err
}

func (r *FolderItemRepositoryImpl) Update(ctx context.Context, id uuid.UUID, item *models.FolderItem) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(item).Error
}

func (r *FolderItemRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.FolderItem{}).Error
}

func (r *FolderItemRepositoryImpl) CountByFolderID(ctx context.Context, folderID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.FolderItem{}).Where("folder_id = ?", folderID).Count(&count).Error
	return count, err
}

func (r *FolderItemRepositoryImpl) UpdateSortOrder(ctx context.Context, id uuid.UUID, sortOrder int) error {
	return r.db.WithContext(ctx).
		Model(&models.FolderItem{}).
		Where("id = ?", id).
		Update("sort_order", sortOrder).Error
}

func (r *FolderItemRepositoryImpl) ReorderItems(ctx context.Context, folderID uuid.UUID, itemOrders map[uuid.UUID]int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for itemID, order := range itemOrders {
			if err := tx.Model(&models.FolderItem{}).
				Where("id = ? AND folder_id = ?", itemID, folderID).
				Update("sort_order", order).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *FolderItemRepositoryImpl) ExistsByFolderIDAndURL(ctx context.Context, folderID uuid.UUID, url string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.FolderItem{}).
		Where("folder_id = ? AND url = ?", folderID, url).
		Count(&count).Error
	return count > 0, err
}

func (r *FolderItemRepositoryImpl) GetFolderIDsByURL(ctx context.Context, userID uuid.UUID, url string) ([]uuid.UUID, error) {
	var folderIDs []uuid.UUID
	err := r.db.WithContext(ctx).
		Model(&models.FolderItem{}).
		Select("folder_items.folder_id").
		Joins("INNER JOIN folders ON folders.id = folder_items.folder_id").
		Where("folders.user_id = ? AND folder_items.url = ?", userID, url).
		Pluck("folder_id", &folderIDs).Error
	return folderIDs, err
}

func (r *FolderItemRepositoryImpl) GetFolderIDsByURLs(ctx context.Context, userID uuid.UUID, urls []string) (map[string][]uuid.UUID, error) {
	type urlFolderPair struct {
		URL      string
		FolderID uuid.UUID
	}

	var pairs []urlFolderPair
	err := r.db.WithContext(ctx).
		Model(&models.FolderItem{}).
		Select("folder_items.url, folder_items.folder_id").
		Joins("INNER JOIN folders ON folders.id = folder_items.folder_id").
		Where("folders.user_id = ? AND folder_items.url IN ?", userID, urls).
		Find(&pairs).Error
	if err != nil {
		return nil, err
	}

	// Group by URL
	result := make(map[string][]uuid.UUID)
	for _, url := range urls {
		result[url] = []uuid.UUID{}
	}
	for _, pair := range pairs {
		result[pair.URL] = append(result[pair.URL], pair.FolderID)
	}

	return result, nil
}
