package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type FavoriteHandler struct {
	favoriteService services.FavoriteService
}

func NewFavoriteHandler(favoriteService services.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
	}
}

func (h *FavoriteHandler) AddFavorite(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.AddFavoriteRequest
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

	result, err := h.favoriteService.AddFavorite(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to add favorite", err)
	}

	return utils.SuccessResponse(c, "Favorite added successfully", result)
}

func (h *FavoriteHandler) GetFavorites(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.GetFavoritesRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid query parameters")
	}

	result, err := h.favoriteService.GetFavorites(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get favorites", err)
	}

	return utils.SuccessResponse(c, "Favorites retrieved successfully", result)
}

func (h *FavoriteHandler) RemoveFavorite(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	favoriteIDStr := c.Params("id")
	favoriteID, err := uuid.Parse(favoriteIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid favorite ID")
	}

	err = h.favoriteService.RemoveFavorite(c.Context(), user.ID, favoriteID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to remove favorite", err)
	}

	return utils.SuccessResponse(c, "Favorite removed successfully", nil)
}

func (h *FavoriteHandler) CheckFavorite(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.CheckFavoriteRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid query parameters")
	}

	result, err := h.favoriteService.CheckFavorite(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check favorite status", err)
	}

	return utils.SuccessResponse(c, "Favorite status retrieved", result)
}

func (h *FavoriteHandler) BatchCheckFavorites(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.BatchCheckFavoritesRequest
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

	result, err := h.favoriteService.BatchCheckFavorites(c.Context(), user.ID, req.ExternalIDs)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to batch check favorites", err)
	}

	return utils.SuccessResponse(c, "Batch favorite check completed", result)
}

func (h *FavoriteHandler) ToggleFavorite(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.AddFavoriteRequest
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

	result, err := h.favoriteService.ToggleFavorite(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to toggle favorite", err)
	}

	message := "Favorite removed"
	if result.IsFavorite {
		message = "Favorite added"
	}

	return utils.SuccessResponse(c, message, result)
}
