package services

import (
	"context"

	"github.com/google/uuid"

	"gofiber-template/domain/dto"
)

type FolderService interface {
	// Folder operations
	CreateFolder(ctx context.Context, userID uuid.UUID, req *dto.CreateFolderRequest) (*dto.FolderResponse, error)
	GetFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID) (*dto.FolderDetailResponse, error)
	GetFolders(ctx context.Context, userID uuid.UUID, req *dto.GetFoldersRequest) (*dto.FolderListResponse, error)
	UpdateFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID, req *dto.UpdateFolderRequest) (*dto.FolderResponse, error)
	DeleteFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID) error

	// Folder item operations
	AddItemToFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID, req *dto.AddFolderItemRequest) (*dto.FolderItemResponse, error)
	GetFolderItems(ctx context.Context, userID uuid.UUID, req *dto.GetFolderItemsRequest) (*dto.FolderItemListResponse, error)
	UpdateFolderItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req *dto.UpdateFolderItemRequest) (*dto.FolderItemResponse, error)
	RemoveItemFromFolder(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error
	ReorderFolderItems(ctx context.Context, userID uuid.UUID, folderID uuid.UUID, req *dto.ReorderFolderItemsRequest) error

	// Folder sharing
	ShareFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID, req *dto.ShareFolderRequest) (*dto.FolderShareResponse, error)
	GetPublicFolder(ctx context.Context, folderID uuid.UUID) (*dto.FolderDetailResponse, error)

	// Check item
	CheckItemInFolders(ctx context.Context, userID uuid.UUID, url string) (*dto.CheckItemInFoldersResponse, error)
}
