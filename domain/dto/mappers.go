package dto

import (
	"encoding/json"

	"gofiber-template/domain/models"
)

func UserToUserResponse(user *models.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Role:      user.Role,
		IsActive:  user.IsActive,
		StudentID: user.StudentID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func CreateUserRequestToUser(req *CreateUserRequest) *models.User {
	return &models.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		StudentID: req.StudentID,
	}
}

func UpdateUserRequestToUser(req *UpdateUserRequest) *models.User {
	return &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Avatar:    req.Avatar,
	}
}

func TaskToTaskResponse(task *models.Task, user *models.User) *TaskResponse {
	if task == nil {
		return nil
	}
	taskResp := &TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		UserID:      task.UserID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
	if user != nil {
		taskResp.User = *UserToUserResponse(user)
	}
	return taskResp
}

func CreateTaskRequestToTask(req *CreateTaskRequest) *models.Task {
	return &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}
}

func UpdateTaskRequestToTask(req *UpdateTaskRequest) *models.Task {
	return &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}
}

func JobToJobResponse(job *models.Job) *JobResponse {
	if job == nil {
		return nil
	}
	return &JobResponse{
		ID:        job.ID,
		Name:      job.Name,
		CronExpr:  job.CronExpr,
		Payload:   job.Payload,
		Status:    job.Status,
		LastRun:   job.LastRun,
		NextRun:   job.NextRun,
		IsActive:  job.IsActive,
		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
	}
}

func CreateJobRequestToJob(req *CreateJobRequest) *models.Job {
	return &models.Job{
		Name:     req.Name,
		CronExpr: req.CronExpr,
		Payload:  req.Payload,
	}
}

func UpdateJobRequestToJob(req *UpdateJobRequest) *models.Job {
	return &models.Job{
		Name:     req.Name,
		CronExpr: req.CronExpr,
		Payload:  req.Payload,
		IsActive: req.IsActive,
	}
}

func FileToFileResponse(file *models.File) *FileResponse {
	if file == nil {
		return nil
	}
	return &FileResponse{
		ID:        file.ID,
		FileName:  file.FileName,
		FileSize:  file.FileSize,
		MimeType:  file.MimeType,
		URL:       file.URL,
		CDNPath:   file.CDNPath,
		UserID:    file.UserID,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
	}
}

// ==================== Folder Mappers ====================

func FolderToFolderResponse(folder *models.Folder) *FolderResponse {
	if folder == nil {
		return nil
	}
	return &FolderResponse{
		ID:            folder.ID,
		Name:          folder.Name,
		Description:   folder.Description,
		CoverImageURL: folder.CoverImageURL,
		IsPublic:      folder.IsPublic,
		ItemCount:     folder.ItemCount,
		CreatedAt:     folder.CreatedAt,
		UpdatedAt:     folder.UpdatedAt,
	}
}

func FolderToFolderDetailResponse(folder *models.Folder) *FolderDetailResponse {
	if folder == nil {
		return nil
	}
	items := make([]FolderItemResponse, 0, len(folder.Items))
	for _, item := range folder.Items {
		items = append(items, *FolderItemToFolderItemResponse(&item))
	}
	return &FolderDetailResponse{
		ID:            folder.ID,
		Name:          folder.Name,
		Description:   folder.Description,
		CoverImageURL: folder.CoverImageURL,
		IsPublic:      folder.IsPublic,
		ItemCount:     folder.ItemCount,
		Items:         items,
		CreatedAt:     folder.CreatedAt,
		UpdatedAt:     folder.UpdatedAt,
	}
}

func CreateFolderRequestToFolder(req *CreateFolderRequest) *models.Folder {
	return &models.Folder{
		Name:          req.Name,
		Description:   req.Description,
		CoverImageURL: req.CoverImageURL,
		IsPublic:      req.IsPublic,
	}
}

func FolderItemToFolderItemResponse(item *models.FolderItem) *FolderItemResponse {
	if item == nil {
		return nil
	}
	var metadata map[string]interface{}
	if item.Metadata != nil {
		_ = json.Unmarshal(item.Metadata, &metadata)
	}
	return &FolderItemResponse{
		ID:           item.ID,
		FolderID:     item.FolderID,
		Type:         item.Type,
		Title:        item.Title,
		URL:          item.URL,
		ThumbnailURL: item.ThumbnailURL,
		Description:  item.Description,
		Metadata:     metadata,
		SortOrder:    item.SortOrder,
		CreatedAt:    item.CreatedAt,
	}
}

func AddFolderItemRequestToFolderItem(req *AddFolderItemRequest) *models.FolderItem {
	var metadataJSON []byte
	if req.Metadata != nil {
		metadataJSON, _ = json.Marshal(req.Metadata)
	}
	return &models.FolderItem{
		Type:         req.Type,
		Title:        req.Title,
		URL:          req.URL,
		ThumbnailURL: req.ThumbnailURL,
		Description:  req.Description,
		Metadata:     metadataJSON,
	}
}

// ==================== Favorite Mappers ====================

func FavoriteToFavoriteResponse(fav *models.Favorite) *FavoriteResponse {
	if fav == nil {
		return nil
	}
	var metadata map[string]interface{}
	if fav.Metadata != nil {
		_ = json.Unmarshal(fav.Metadata, &metadata)
	}
	return &FavoriteResponse{
		ID:           fav.ID,
		Type:         fav.Type,
		ExternalID:   fav.ExternalID,
		Title:        fav.Title,
		URL:          fav.URL,
		ThumbnailURL: fav.ThumbnailURL,
		Rating:       fav.Rating,
		ReviewCount:  fav.ReviewCount,
		Address:      fav.Address,
		Metadata:     metadata,
		CreatedAt:    fav.CreatedAt,
	}
}

func AddFavoriteRequestToFavorite(req *AddFavoriteRequest) *models.Favorite {
	var metadataJSON []byte
	if req.Metadata != nil {
		metadataJSON, _ = json.Marshal(req.Metadata)
	}
	return &models.Favorite{
		Type:         req.Type,
		ExternalID:   req.ExternalID,
		Title:        req.Title,
		URL:          req.URL,
		ThumbnailURL: req.ThumbnailURL,
		Rating:       req.Rating,
		ReviewCount:  req.ReviewCount,
		Address:      req.Address,
		Metadata:     metadataJSON,
	}
}

// ==================== Search History Mappers ====================

func SearchHistoryToSearchHistoryResponse(sh *models.SearchHistory) *SearchHistoryResponse {
	if sh == nil {
		return nil
	}
	return &SearchHistoryResponse{
		ID:          sh.ID,
		Query:       sh.Query,
		SearchType:  sh.SearchType,
		ResultCount: sh.ResultCount,
		CreatedAt:   sh.CreatedAt,
	}
}

// ==================== AI Chat Mappers ====================

func AIChatSessionToResponse(session *models.AIChatSession) *AIChatSessionResponse {
	if session == nil {
		return nil
	}
	return &AIChatSessionResponse{
		ID:           session.ID,
		Title:        session.Title,
		InitialQuery: session.InitialQuery,
		CreatedAt:    session.CreatedAt,
		UpdatedAt:    session.UpdatedAt,
	}
}

func AIChatSessionToDetailResponse(session *models.AIChatSession) *AIChatSessionDetailResponse {
	if session == nil {
		return nil
	}
	messages := make([]AIChatMessageResponse, 0, len(session.Messages))
	for _, msg := range session.Messages {
		messages = append(messages, *AIChatMessageToResponse(&msg))
	}
	return &AIChatSessionDetailResponse{
		ID:           session.ID,
		Title:        session.Title,
		InitialQuery: session.InitialQuery,
		Messages:     messages,
		CreatedAt:    session.CreatedAt,
		UpdatedAt:    session.UpdatedAt,
	}
}

func AIChatMessageToResponse(msg *models.AIChatMessage) *AIChatMessageResponse {
	if msg == nil {
		return nil
	}
	var sources []MessageSource
	if msg.Sources != nil {
		_ = json.Unmarshal(msg.Sources, &sources)
	}
	return &AIChatMessageResponse{
		ID:        msg.ID,
		SessionID: msg.SessionID,
		Role:      msg.Role,
		Content:   msg.Content,
		Sources:   sources,
		CreatedAt: msg.CreatedAt,
	}
}