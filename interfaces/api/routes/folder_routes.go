package routes

import (
	"github.com/gofiber/fiber/v2"

	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupFolderRoutes(api fiber.Router, h *handlers.Handlers) {
	folders := api.Group("/folders")

	// Public folder access
	folders.Get("/public/:id", h.FolderHandler.GetPublicFolder)

	// Protected folder endpoints
	protected := folders.Group("/")
	protected.Use(middleware.Protected())
	protected.Post("/", h.FolderHandler.CreateFolder)
	protected.Get("/", h.FolderHandler.GetFolders)
	protected.Get("/items/check", h.FolderHandler.CheckItemInFolders)
	protected.Post("/items/check/batch", h.FolderHandler.BatchCheckItemsInFolders)
	protected.Get("/:id", h.FolderHandler.GetFolder)
	protected.Put("/:id", h.FolderHandler.UpdateFolder)
	protected.Delete("/:id", h.FolderHandler.DeleteFolder)
	protected.Post("/:id/share", h.FolderHandler.ShareFolder)

	// Folder items
	protected.Post("/:id/items", h.FolderHandler.AddItemToFolder)
	protected.Post("/:id/items/upload", h.FolderHandler.UploadItemToFolder)
	protected.Get("/:id/items", h.FolderHandler.GetFolderItems)
	protected.Put("/:id/items/reorder", h.FolderHandler.ReorderFolderItems)
	protected.Put("/items/:itemId", h.FolderHandler.UpdateFolderItem)
	protected.Delete("/items/:itemId", h.FolderHandler.RemoveItemFromFolder)
}
