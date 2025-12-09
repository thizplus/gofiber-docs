package handlers

import (
	"github.com/gofiber/fiber/v2"

	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/config"
	"gofiber-template/pkg/utils"
)

type UtilityHandler struct {
	utilityService services.UtilityService
	config         *config.Config
}

func NewUtilityHandler(utilityService services.UtilityService, cfg *config.Config) *UtilityHandler {
	return &UtilityHandler{
		utilityService: utilityService,
		config:         cfg,
	}
}

func (h *UtilityHandler) Translate(c *fiber.Ctx) error {
	var req dto.TranslateRequest
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

	result, err := h.utilityService.Translate(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Translation failed", err)
	}

	return utils.SuccessResponse(c, "Translation completed", result)
}

func (h *UtilityHandler) DetectLanguage(c *fiber.Ctx) error {
	var req dto.DetectLanguageRequest
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

	result, err := h.utilityService.DetectLanguage(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Language detection failed", err)
	}

	return utils.SuccessResponse(c, "Language detected", result)
}

func (h *UtilityHandler) GenerateQRCode(c *fiber.Ctx) error {
	var req dto.GenerateQRRequest
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

	result, err := h.utilityService.GenerateQRCode(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "QR code generation failed", err)
	}

	return utils.SuccessResponse(c, "QR code generated", result)
}

func (h *UtilityHandler) CalculateDistance(c *fiber.Ctx) error {
	var req dto.CalculateDistanceRequest
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

	result, err := h.utilityService.CalculateDistance(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Distance calculation failed", err)
	}

	return utils.SuccessResponse(c, "Distance calculated", result)
}

func (h *UtilityHandler) HealthCheck(c *fiber.Ctx) error {
	result, err := h.utilityService.HealthCheck(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Health check failed", err)
	}

	return utils.SuccessResponse(c, "System healthy", result)
}

// GetPublicConfig returns public configuration for frontend
func (h *UtilityHandler) GetPublicConfig(c *fiber.Ctx) error {
	publicConfig := fiber.Map{
		"googleMapsApiKey": h.config.Google.MapsAPIKey,
		"appName":          h.config.App.Name,
	}

	return utils.SuccessResponse(c, "Public config retrieved", publicConfig)
}
