package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type SearchHandler struct {
	searchService services.SearchService
}

func NewSearchHandler(searchService services.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

func (h *SearchHandler) Search(c *fiber.Ctx) error {
	var req dto.SearchRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid query parameters")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	userID := getUserIDFromContext(c)

	result, err := h.searchService.Search(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Search failed", err)
	}

	return utils.SuccessResponse(c, "Search completed", result)
}

func (h *SearchHandler) SearchWebsites(c *fiber.Ctx) error {
	var req dto.SearchRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid query parameters")
	}

	userID := getUserIDFromContext(c)

	result, err := h.searchService.SearchWebsites(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Website search failed", err)
	}

	return utils.SuccessResponse(c, "Website search completed", result)
}

func (h *SearchHandler) SearchImages(c *fiber.Ctx) error {
	var req dto.ImageSearchRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid query parameters")
	}

	userID := getUserIDFromContext(c)

	result, err := h.searchService.SearchImages(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Image search failed", err)
	}

	return utils.SuccessResponse(c, "Image search completed", result)
}

func (h *SearchHandler) SearchVideos(c *fiber.Ctx) error {
	var req dto.VideoSearchRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid query parameters")
	}

	userID := getUserIDFromContext(c)

	result, err := h.searchService.SearchVideos(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Video search failed", err)
	}

	return utils.SuccessResponse(c, "Video search completed", result)
}

func (h *SearchHandler) GetVideoDetails(c *fiber.Ctx) error {
	videoID := c.Params("videoId")
	if videoID == "" {
		return utils.ValidationErrorResponse(c, "Video ID is required")
	}

	result, err := h.searchService.GetVideoDetails(c.Context(), videoID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get video details", err)
	}

	return utils.SuccessResponse(c, "Video details retrieved", result)
}

func (h *SearchHandler) SearchPlaces(c *fiber.Ctx) error {
	var req dto.PlaceSearchRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid query parameters")
	}

	userID := getUserIDFromContext(c)

	result, err := h.searchService.SearchPlaces(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Place search failed", err)
	}

	return utils.SuccessResponse(c, "Place search completed", result)
}

func (h *SearchHandler) GetPlaceDetails(c *fiber.Ctx) error {
	placeID := c.Params("placeId")
	if placeID == "" {
		return utils.ValidationErrorResponse(c, "Place ID is required")
	}

	userLat := c.QueryFloat("lat", 0)
	userLng := c.QueryFloat("lng", 0)

	result, err := h.searchService.GetPlaceDetails(c.Context(), placeID, userLat, userLng)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get place details", err)
	}

	return utils.SuccessResponse(c, "Place details retrieved", result)
}

func (h *SearchHandler) SearchNearbyPlaces(c *fiber.Ctx) error {
	var req dto.NearbyPlacesRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid query parameters")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	result, err := h.searchService.SearchNearbyPlaces(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Nearby places search failed", err)
	}

	return utils.SuccessResponse(c, "Nearby places search completed", result)
}

func (h *SearchHandler) GetSearchHistory(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.GetSearchHistoryRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid query parameters")
	}

	result, err := h.searchService.GetSearchHistory(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get search history", err)
	}

	return utils.SuccessResponse(c, "Search history retrieved", result)
}

func (h *SearchHandler) ClearSearchHistory(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.ClearSearchHistoryRequest
	if err := c.BodyParser(&req); err != nil {
		req = dto.ClearSearchHistoryRequest{}
	}

	err = h.searchService.ClearSearchHistory(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to clear search history", err)
	}

	return utils.SuccessResponse(c, "Search history cleared", nil)
}

func (h *SearchHandler) DeleteSearchHistoryItem(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	historyIDStr := c.Params("id")
	historyID, err := uuid.Parse(historyIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid history ID")
	}

	err = h.searchService.DeleteSearchHistoryItem(c.Context(), user.ID, historyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete history item", err)
	}

	return utils.SuccessResponse(c, "History item deleted", nil)
}

// Helper function to get user ID from context (may be nil for unauthenticated users)
func getUserIDFromContext(c *fiber.Ctx) uuid.UUID {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return uuid.Nil
	}
	return user.ID
}
