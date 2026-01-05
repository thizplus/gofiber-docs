package serviceimpl

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/infrastructure/storage"
)

type FolderServiceImpl struct {
	folderRepo     repositories.FolderRepository
	folderItemRepo repositories.FolderItemRepository
	r2Storage      storage.R2Storage
}

func NewFolderService(
	folderRepo repositories.FolderRepository,
	folderItemRepo repositories.FolderItemRepository,
	r2Storage storage.R2Storage,
) services.FolderService {
	return &FolderServiceImpl{
		folderRepo:     folderRepo,
		folderItemRepo: folderItemRepo,
		r2Storage:      r2Storage,
	}
}

func (s *FolderServiceImpl) CreateFolder(ctx context.Context, userID uuid.UUID, req *dto.CreateFolderRequest) (*dto.FolderResponse, error) {
	folder := dto.CreateFolderRequestToFolder(req)
	folder.UserID = userID
	folder.CreatedAt = time.Now()
	folder.UpdatedAt = time.Now()

	if err := s.folderRepo.Create(ctx, folder); err != nil {
		return nil, err
	}

	return dto.FolderToFolderResponse(folder), nil
}

func (s *FolderServiceImpl) GetFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID) (*dto.FolderDetailResponse, error) {
	folder, err := s.folderRepo.GetByIDWithItems(ctx, folderID)
	if err != nil {
		return nil, errors.New("folder not found")
	}

	// Check ownership or public access
	if folder.UserID != userID && !folder.IsPublic {
		return nil, errors.New("unauthorized")
	}

	return dto.FolderToFolderDetailResponse(folder), nil
}

func (s *FolderServiceImpl) GetFolders(ctx context.Context, userID uuid.UUID, req *dto.GetFoldersRequest) (*dto.FolderListResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	folders, err := s.folderRepo.GetByUserID(ctx, userID, offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	total, err := s.folderRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var folderResponses []dto.FolderResponse
	for _, f := range folders {
		folderResponses = append(folderResponses, *dto.FolderToFolderResponse(f))
	}

	return &dto.FolderListResponse{
		Folders: folderResponses,
		Meta: dto.PaginationMeta{
			Total:  total,
			Offset: offset,
			Limit:  req.PageSize,
		},
	}, nil
}

func (s *FolderServiceImpl) UpdateFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID, req *dto.UpdateFolderRequest) (*dto.FolderResponse, error) {
	folder, err := s.folderRepo.GetByID(ctx, folderID)
	if err != nil {
		return nil, errors.New("folder not found")
	}

	if folder.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	if req.Name != "" {
		folder.Name = req.Name
	}
	if req.Description != "" {
		folder.Description = req.Description
	}
	if req.CoverImageURL != "" {
		folder.CoverImageURL = req.CoverImageURL
	}
	if req.IsPublic != nil {
		folder.IsPublic = *req.IsPublic
	}
	folder.UpdatedAt = time.Now()

	if err := s.folderRepo.Update(ctx, folderID, folder); err != nil {
		return nil, err
	}

	return dto.FolderToFolderResponse(folder), nil
}

func (s *FolderServiceImpl) DeleteFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID) error {
	folder, err := s.folderRepo.GetByID(ctx, folderID)
	if err != nil {
		return errors.New("folder not found")
	}

	if folder.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.folderRepo.Delete(ctx, folderID)
}

func (s *FolderServiceImpl) AddItemToFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID, req *dto.AddFolderItemRequest) (*dto.FolderItemResponse, error) {
	folder, err := s.folderRepo.GetByID(ctx, folderID)
	if err != nil {
		return nil, errors.New("folder not found")
	}

	if folder.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Check if item already exists
	exists, err := s.folderItemRepo.ExistsByFolderIDAndURL(ctx, folderID, req.URL)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("item already exists in folder")
	}

	item := dto.AddFolderItemRequestToFolderItem(req)
	item.FolderID = folderID
	item.CreatedAt = time.Now()

	if err := s.folderItemRepo.Create(ctx, item); err != nil {
		return nil, err
	}

	// Increment folder item count
	_ = s.folderRepo.IncrementItemCount(ctx, folderID)

	return dto.FolderItemToFolderItemResponse(item), nil
}

func (s *FolderServiceImpl) GetFolderItems(ctx context.Context, userID uuid.UUID, req *dto.GetFolderItemsRequest) (*dto.FolderItemListResponse, error) {
	folder, err := s.folderRepo.GetByID(ctx, req.FolderID)
	if err != nil {
		return nil, errors.New("folder not found")
	}

	if folder.UserID != userID && !folder.IsPublic {
		return nil, errors.New("unauthorized")
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	var items []*models.FolderItem
	if req.Type != "" {
		items, err = s.folderItemRepo.GetByFolderIDAndType(ctx, req.FolderID, req.Type, offset, req.PageSize)
	} else {
		items, err = s.folderItemRepo.GetByFolderID(ctx, req.FolderID, offset, req.PageSize)
	}
	if err != nil {
		return nil, err
	}

	total, err := s.folderItemRepo.CountByFolderID(ctx, req.FolderID)
	if err != nil {
		return nil, err
	}

	var itemResponses []dto.FolderItemResponse
	for _, i := range items {
		itemResponses = append(itemResponses, *dto.FolderItemToFolderItemResponse(i))
	}

	return &dto.FolderItemListResponse{
		Items: itemResponses,
		Meta: dto.PaginationMeta{
			Total:  total,
			Offset: offset,
			Limit:  req.PageSize,
		},
	}, nil
}

func (s *FolderServiceImpl) UpdateFolderItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req *dto.UpdateFolderItemRequest) (*dto.FolderItemResponse, error) {
	item, err := s.folderItemRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	folder, err := s.folderRepo.GetByID(ctx, item.FolderID)
	if err != nil {
		return nil, errors.New("folder not found")
	}

	if folder.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	if req.Title != "" {
		item.Title = req.Title
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.SortOrder != nil {
		item.SortOrder = *req.SortOrder
	}

	if err := s.folderItemRepo.Update(ctx, itemID, item); err != nil {
		return nil, err
	}

	return dto.FolderItemToFolderItemResponse(item), nil
}

func (s *FolderServiceImpl) RemoveItemFromFolder(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error {
	item, err := s.folderItemRepo.GetByID(ctx, itemID)
	if err != nil {
		return errors.New("item not found")
	}

	folder, err := s.folderRepo.GetByID(ctx, item.FolderID)
	if err != nil {
		return errors.New("folder not found")
	}

	if folder.UserID != userID {
		return errors.New("unauthorized")
	}

	if err := s.folderItemRepo.Delete(ctx, itemID); err != nil {
		return err
	}

	// Decrement folder item count
	_ = s.folderRepo.DecrementItemCount(ctx, item.FolderID)

	return nil
}

func (s *FolderServiceImpl) ReorderFolderItems(ctx context.Context, userID uuid.UUID, folderID uuid.UUID, req *dto.ReorderFolderItemsRequest) error {
	folder, err := s.folderRepo.GetByID(ctx, folderID)
	if err != nil {
		return errors.New("folder not found")
	}

	if folder.UserID != userID {
		return errors.New("unauthorized")
	}

	itemOrders := make(map[uuid.UUID]int)
	for _, order := range req.ItemOrders {
		itemOrders[order.ItemID] = order.SortOrder
	}

	return s.folderItemRepo.ReorderItems(ctx, folderID, itemOrders)
}

func (s *FolderServiceImpl) ShareFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID, req *dto.ShareFolderRequest) (*dto.FolderShareResponse, error) {
	folder, err := s.folderRepo.GetByID(ctx, folderID)
	if err != nil {
		return nil, errors.New("folder not found")
	}

	if folder.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	if req.IsPublic != nil {
		folder.IsPublic = *req.IsPublic
	}
	folder.UpdatedAt = time.Now()

	if err := s.folderRepo.Update(ctx, folderID, folder); err != nil {
		return nil, err
	}

	response := &dto.FolderShareResponse{
		FolderID: folderID,
		IsPublic: folder.IsPublic,
	}

	if folder.IsPublic {
		response.ShareURL = "/folders/public/" + folderID.String()
	}

	return response, nil
}

func (s *FolderServiceImpl) GetPublicFolder(ctx context.Context, folderID uuid.UUID) (*dto.FolderDetailResponse, error) {
	folder, err := s.folderRepo.GetByIDWithItems(ctx, folderID)
	if err != nil {
		return nil, errors.New("folder not found")
	}

	if !folder.IsPublic {
		return nil, errors.New("folder is not public")
	}

	return dto.FolderToFolderDetailResponse(folder), nil
}

func (s *FolderServiceImpl) CheckItemInFolders(ctx context.Context, userID uuid.UUID, url string) (*dto.CheckItemInFoldersResponse, error) {
	folderIDs, err := s.folderItemRepo.GetFolderIDsByURL(ctx, userID, url)
	if err != nil {
		return nil, err
	}

	response := &dto.CheckItemInFoldersResponse{
		IsSaved:   len(folderIDs) > 0,
		FolderIDs: folderIDs,
		Folders:   make([]dto.FolderSummaryResponse, 0),
	}

	// Get folder names
	for _, folderID := range folderIDs {
		folder, err := s.folderRepo.GetByID(ctx, folderID)
		if err == nil {
			response.Folders = append(response.Folders, dto.FolderSummaryResponse{
				ID:   folder.ID,
				Name: folder.Name,
			})
		}
	}

	return response, nil
}

func (s *FolderServiceImpl) BatchCheckItemsInFolders(ctx context.Context, userID uuid.UUID, urls []string) (*dto.BatchCheckItemsResponse, error) {
	urlFolderMap, err := s.folderItemRepo.GetFolderIDsByURLs(ctx, userID, urls)
	if err != nil {
		return nil, err
	}

	// Cache folders to avoid duplicate queries
	folderCache := make(map[uuid.UUID]*dto.FolderSummaryResponse)

	response := &dto.BatchCheckItemsResponse{
		Items: make(map[string]dto.CheckItemInFoldersResponse),
	}

	for url, folderIDs := range urlFolderMap {
		itemResponse := dto.CheckItemInFoldersResponse{
			IsSaved:   len(folderIDs) > 0,
			FolderIDs: folderIDs,
			Folders:   make([]dto.FolderSummaryResponse, 0),
		}

		// Get folder names (with caching)
		for _, folderID := range folderIDs {
			if cached, ok := folderCache[folderID]; ok {
				itemResponse.Folders = append(itemResponse.Folders, *cached)
			} else {
				folder, err := s.folderRepo.GetByID(ctx, folderID)
				if err == nil {
					summary := &dto.FolderSummaryResponse{
						ID:   folder.ID,
						Name: folder.Name,
					}
					folderCache[folderID] = summary
					itemResponse.Folders = append(itemResponse.Folders, *summary)
				}
			}
		}

		response.Items[url] = itemResponse
	}

	return response, nil
}

// Upload file size limits
const (
	maxImageSize = 10 * 1024 * 1024  // 10MB
	maxPDFSize   = 20 * 1024 * 1024  // 20MB
	maxVideoSize = 100 * 1024 * 1024 // 100MB
)

// Allowed file extensions
var (
	imageExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	pdfExtensions   = []string{".pdf"}
	videoExtensions = []string{".mp4", ".mov", ".webm"}
)

func (s *FolderServiceImpl) UploadItemToFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID, file *multipart.FileHeader) (*dto.UploadItemResponse, error) {
	// Check folder ownership
	folder, err := s.folderRepo.GetByID(ctx, folderID)
	if err != nil {
		return nil, errors.New("folder not found")
	}

	if folder.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Check if R2 storage is available
	if s.r2Storage == nil {
		return nil, errors.New("storage not configured")
	}

	// Detect file type and validate
	ext := strings.ToLower(filepath.Ext(file.Filename))
	fileType, maxSize := detectFileTypeAndMaxSize(ext)
	if fileType == "" {
		return nil, errors.New("unsupported file type")
	}

	if file.Size > maxSize {
		return nil, fmt.Errorf("file too large. Max size for %s is %dMB", fileType, maxSize/(1024*1024))
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return nil, errors.New("failed to read file")
	}
	defer src.Close()

	// Generate unique filename and path
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	storagePath := fmt.Sprintf("folders/%s/%s/%s", folderID.String(), getStorageFolder(fileType), uniqueFilename)

	// Get content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = getContentType(ext)
	}

	// Upload to R2
	fileURL, err := s.r2Storage.UploadFile(src, storagePath, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Create folder item
	item := &models.FolderItem{
		FolderID:     folderID,
		Type:         fileType,
		Title:        sanitizeFilename(file.Filename),
		URL:          fileURL,
		ThumbnailURL: fileURL, // For images, use the same URL as thumbnail
		Description:  fmt.Sprintf("Uploaded %s file", fileType),
		CreatedAt:    time.Now(),
	}

	if err := s.folderItemRepo.Create(ctx, item); err != nil {
		// Try to delete uploaded file on failure
		_ = s.r2Storage.DeleteFile(storagePath)
		return nil, errors.New("failed to save file record")
	}

	// Increment folder item count
	_ = s.folderRepo.IncrementItemCount(ctx, folderID)

	return &dto.UploadItemResponse{
		Item:     *dto.FolderItemToFolderItemResponse(item),
		FileURL:  fileURL,
		FileSize: file.Size,
		MimeType: contentType,
	}, nil
}

func detectFileTypeAndMaxSize(ext string) (string, int64) {
	for _, e := range imageExtensions {
		if e == ext {
			return "image", maxImageSize
		}
	}
	for _, e := range pdfExtensions {
		if e == ext {
			return "pdf", maxPDFSize
		}
	}
	for _, e := range videoExtensions {
		if e == ext {
			return "video", maxVideoSize
		}
	}
	return "", 0
}

func getStorageFolder(fileType string) string {
	switch fileType {
	case "image":
		return "images"
	case "pdf":
		return "documents"
	case "video":
		return "videos"
	default:
		return "uploads"
	}
}

func getContentType(ext string) string {
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".pdf":  "application/pdf",
		".mp4":  "video/mp4",
		".mov":  "video/quicktime",
		".webm": "video/webm",
	}
	if ct, ok := contentTypes[ext]; ok {
		return ct
	}
	return "application/octet-stream"
}

func sanitizeFilename(filename string) string {
	// Remove path and keep only filename
	name := filepath.Base(filename)
	// Remove extension for title
	ext := filepath.Ext(name)
	if ext != "" {
		name = name[:len(name)-len(ext)]
	}
	return name
}
