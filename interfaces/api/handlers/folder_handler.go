package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type FolderHandler struct {
	folderService services.FolderService
}

func NewFolderHandler(folderService services.FolderService) *FolderHandler {
	return &FolderHandler{
		folderService: folderService,
	}
}

func (h *FolderHandler) CreateFolder(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.CreateFolderRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	result, err := h.folderService.CreateFolder(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create folder", err)
	}

	return utils.SuccessResponse(c, "Folder created successfully", result)
}

func (h *FolderHandler) GetFolder(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	folderIDStr := c.Params("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid folder ID")
	}

	result, err := h.folderService.GetFolder(c.Context(), user.ID, folderID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Folder not found", err)
	}

	return utils.SuccessResponse(c, "Folder retrieved successfully", result)
}

func (h *FolderHandler) GetFolders(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.GetFoldersRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid query parameters")
	}

	result, err := h.folderService.GetFolders(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get folders", err)
	}

	return utils.SuccessResponse(c, "Folders retrieved successfully", result)
}

func (h *FolderHandler) UpdateFolder(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	folderIDStr := c.Params("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid folder ID")
	}

	var req dto.UpdateFolderRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	result, err := h.folderService.UpdateFolder(c.Context(), user.ID, folderID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update folder", err)
	}

	return utils.SuccessResponse(c, "Folder updated successfully", result)
}

func (h *FolderHandler) DeleteFolder(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	folderIDStr := c.Params("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid folder ID")
	}

	err = h.folderService.DeleteFolder(c.Context(), user.ID, folderID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete folder", err)
	}

	return utils.SuccessResponse(c, "Folder deleted successfully", nil)
}

func (h *FolderHandler) AddItemToFolder(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	folderIDStr := c.Params("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid folder ID")
	}

	var req dto.AddFolderItemRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	result, err := h.folderService.AddItemToFolder(c.Context(), user.ID, folderID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to add item to folder", err)
	}

	return utils.SuccessResponse(c, "Item added to folder successfully", result)
}

func (h *FolderHandler) GetFolderItems(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	folderIDStr := c.Params("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid folder ID")
	}

	var req dto.GetFolderItemsRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid query parameters")
	}
	req.FolderID = folderID

	result, err := h.folderService.GetFolderItems(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get folder items", err)
	}

	return utils.SuccessResponse(c, "Folder items retrieved successfully", result)
}

func (h *FolderHandler) UpdateFolderItem(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	itemIDStr := c.Params("itemId")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid item ID")
	}

	var req dto.UpdateFolderItemRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	result, err := h.folderService.UpdateFolderItem(c.Context(), user.ID, itemID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update folder item", err)
	}

	return utils.SuccessResponse(c, "Folder item updated successfully", result)
}

func (h *FolderHandler) RemoveItemFromFolder(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	itemIDStr := c.Params("itemId")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid item ID")
	}

	err = h.folderService.RemoveItemFromFolder(c.Context(), user.ID, itemID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to remove item from folder", err)
	}

	return utils.SuccessResponse(c, "Item removed from folder successfully", nil)
}

func (h *FolderHandler) ReorderFolderItems(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	folderIDStr := c.Params("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid folder ID")
	}

	var req dto.ReorderFolderItemsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	err = h.folderService.ReorderFolderItems(c.Context(), user.ID, folderID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to reorder folder items", err)
	}

	return utils.SuccessResponse(c, "Folder items reordered successfully", nil)
}

func (h *FolderHandler) ShareFolder(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	folderIDStr := c.Params("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid folder ID")
	}

	var req dto.ShareFolderRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	result, err := h.folderService.ShareFolder(c.Context(), user.ID, folderID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to share folder", err)
	}

	return utils.SuccessResponse(c, "Folder sharing updated successfully", result)
}

func (h *FolderHandler) GetPublicFolder(c *fiber.Ctx) error {
	folderIDStr := c.Params("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid folder ID")
	}

	result, err := h.folderService.GetPublicFolder(c.Context(), folderID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Folder not found or not public", err)
	}

	return utils.SuccessResponse(c, "Public folder retrieved successfully", result)
}

func (h *FolderHandler) CheckItemInFolders(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.CheckItemInFoldersRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid query parameters")
	}

	if req.URL == "" {
		return utils.ValidationErrorResponse(c, "URL is required")
	}

	result, err := h.folderService.CheckItemInFolders(c.Context(), user.ID, req.URL)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check item", err)
	}

	return utils.SuccessResponse(c, "Item check completed", result)
}
