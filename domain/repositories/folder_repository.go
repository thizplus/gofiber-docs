package repositories

import (
	"context"

	"github.com/google/uuid"

	"gofiber-template/domain/models"
)

type FolderRepository interface {
	// Folder operations
	Create(ctx context.Context, folder *models.Folder) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Folder, error)
	GetByIDWithItems(ctx context.Context, id uuid.UUID) (*models.Folder, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Folder, error)
	GetPublicFolders(ctx context.Context, offset, limit int) ([]*models.Folder, error)
	Update(ctx context.Context, id uuid.UUID, folder *models.Folder) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	IncrementItemCount(ctx context.Context, id uuid.UUID) error
	DecrementItemCount(ctx context.Context, id uuid.UUID) error
}

type FolderItemRepository interface {
	// Folder item operations
	Create(ctx context.Context, item *models.FolderItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.FolderItem, error)
	GetByFolderID(ctx context.Context, folderID uuid.UUID, offset, limit int) ([]*models.FolderItem, error)
	GetByFolderIDAndType(ctx context.Context, folderID uuid.UUID, itemType string, offset, limit int) ([]*models.FolderItem, error)
	Update(ctx context.Context, id uuid.UUID, item *models.FolderItem) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByFolderID(ctx context.Context, folderID uuid.UUID) (int64, error)
	UpdateSortOrder(ctx context.Context, id uuid.UUID, sortOrder int) error
	ReorderItems(ctx context.Context, folderID uuid.UUID, itemOrders map[uuid.UUID]int) error
	ExistsByFolderIDAndURL(ctx context.Context, folderID uuid.UUID, url string) (bool, error)
	GetFolderIDsByURL(ctx context.Context, userID uuid.UUID, url string) ([]uuid.UUID, error)
}
