package dto

import (
	"time"

	"github.com/google/uuid"
)

// ==================== Folder DTOs ====================

type CreateFolderRequest struct {
	Name          string `json:"name" validate:"required,min=1,max=255"`
	Description   string `json:"description" validate:"omitempty,max=1000"`
	CoverImageURL string `json:"coverImageUrl" validate:"omitempty,url,max=2000"`
	IsPublic      bool   `json:"isPublic"`
}

type UpdateFolderRequest struct {
	Name          string `json:"name" validate:"omitempty,min=1,max=255"`
	Description   string `json:"description" validate:"omitempty,max=1000"`
	CoverImageURL string `json:"coverImageUrl" validate:"omitempty,url,max=2000"`
	IsPublic      *bool  `json:"isPublic"`
}

type FolderResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	CoverImageURL string    `json:"coverImageUrl,omitempty"`
	IsPublic      bool      `json:"isPublic"`
	ItemCount     int       `json:"itemCount"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type FolderDetailResponse struct {
	ID            uuid.UUID            `json:"id"`
	Name          string               `json:"name"`
	Description   string               `json:"description"`
	CoverImageURL string               `json:"coverImageUrl,omitempty"`
	IsPublic      bool                 `json:"isPublic"`
	ItemCount     int                  `json:"itemCount"`
	Items         []FolderItemResponse `json:"items"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
}

type FolderListResponse struct {
	Folders []FolderResponse `json:"folders"`
	Meta    PaginationMeta   `json:"meta"`
}

type GetFoldersRequest struct {
	Page     int  `json:"page" query:"page" validate:"omitempty,min=1"`
	PageSize int  `json:"pageSize" query:"pageSize" validate:"omitempty,min=1,max=50"`
	IsPublic *bool `json:"isPublic" query:"isPublic"`
}

// ==================== Folder Item DTOs ====================

type AddFolderItemRequest struct {
	Type         string                 `json:"type" validate:"required,oneof=place website image video link"`
	Title        string                 `json:"title" validate:"required,min=1,max=255"`
	URL          string                 `json:"url" validate:"required,url,max=2000"`
	ThumbnailURL string                 `json:"thumbnailUrl" validate:"omitempty,url,max=2000"`
	Description  string                 `json:"description" validate:"omitempty,max=1000"`
	Metadata     map[string]interface{} `json:"metadata" validate:"omitempty"`
}

type UpdateFolderItemRequest struct {
	Title       string `json:"title" validate:"omitempty,min=1,max=255"`
	Description string `json:"description" validate:"omitempty,max=1000"`
	SortOrder   *int   `json:"sortOrder" validate:"omitempty,min=0"`
}

type FolderItemResponse struct {
	ID           uuid.UUID              `json:"id"`
	FolderID     uuid.UUID              `json:"folderId"`
	Type         string                 `json:"type"`
	Title        string                 `json:"title"`
	URL          string                 `json:"url"`
	ThumbnailURL string                 `json:"thumbnailUrl,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	SortOrder    int                    `json:"sortOrder"`
	CreatedAt    time.Time              `json:"createdAt"`
}

type FolderItemListResponse struct {
	Items []FolderItemResponse `json:"items"`
	Meta  PaginationMeta       `json:"meta"`
}

type GetFolderItemsRequest struct {
	FolderID uuid.UUID `json:"folderId" param:"folderId" validate:"required"`
	Type     string    `json:"type" query:"type" validate:"omitempty,oneof=place website image video link"`
	Page     int       `json:"page" query:"page" validate:"omitempty,min=1"`
	PageSize int       `json:"pageSize" query:"pageSize" validate:"omitempty,min=1,max=50"`
}

type ReorderFolderItemsRequest struct {
	ItemOrders []ItemOrder `json:"itemOrders" validate:"required,min=1,dive"`
}

type ItemOrder struct {
	ItemID    uuid.UUID `json:"itemId" validate:"required"`
	SortOrder int       `json:"sortOrder" validate:"min=0"`
}

// ==================== Folder Sharing DTOs ====================

type ShareFolderRequest struct {
	IsPublic *bool `json:"isPublic" validate:"required"`
}

type FolderShareResponse struct {
	FolderID uuid.UUID `json:"folderId"`
	IsPublic bool      `json:"isPublic"`
	ShareURL string    `json:"shareUrl,omitempty"`
}

// ==================== Check Item DTOs ====================

type CheckItemInFoldersRequest struct {
	URL string `json:"url" query:"url" validate:"required,url"`
}

type CheckItemInFoldersResponse struct {
	IsSaved   bool                    `json:"isSaved"`
	FolderIDs []uuid.UUID             `json:"folderIds"`
	Folders   []FolderSummaryResponse `json:"folders"`
}

type FolderSummaryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// ==================== Batch Check Item DTOs ====================

type BatchCheckItemsRequest struct {
	URLs []string `json:"urls" validate:"required,min=1,max=50,dive,url"`
}

type BatchCheckItemsResponse struct {
	Items map[string]CheckItemInFoldersResponse `json:"items"`
}

// ==================== Upload Item DTOs ====================

type UploadItemResponse struct {
	Item     FolderItemResponse `json:"item"`
	FileURL  string             `json:"fileUrl"`
	FileSize int64              `json:"fileSize"`
	MimeType string             `json:"mimeType"`
}
