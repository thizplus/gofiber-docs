package handlers

import (
	"gofiber-template/domain/services"
	"gofiber-template/pkg/config"
)

// Services contains all the services needed for handlers
type Services struct {
	UserService     services.UserService
	TaskService     services.TaskService
	FileService     services.FileService
	JobService      services.JobService
	SearchService   services.SearchService
	AIService       services.AIService
	FolderService   services.FolderService
	FavoriteService services.FavoriteService
	UtilityService  services.UtilityService
}

// Handlers contains all HTTP handlers
type Handlers struct {
	UserHandler     *UserHandler
	TaskHandler     *TaskHandler
	FileHandler     *FileHandler
	JobHandler      *JobHandler
	SearchHandler   *SearchHandler
	AIHandler       *AIHandler
	FolderHandler   *FolderHandler
	FavoriteHandler *FavoriteHandler
	UtilityHandler  *UtilityHandler
}

// NewHandlers creates a new instance of Handlers with all dependencies
func NewHandlers(services *Services, cfg *config.Config) *Handlers {
	return &Handlers{
		UserHandler:     NewUserHandler(services.UserService),
		TaskHandler:     NewTaskHandler(services.TaskService),
		FileHandler:     NewFileHandler(services.FileService),
		JobHandler:      NewJobHandler(services.JobService),
		SearchHandler:   NewSearchHandler(services.SearchService),
		AIHandler:       NewAIHandler(services.AIService),
		FolderHandler:   NewFolderHandler(services.FolderService),
		FavoriteHandler: NewFavoriteHandler(services.FavoriteService),
		UtilityHandler:  NewUtilityHandler(services.UtilityService, cfg),
	}
}