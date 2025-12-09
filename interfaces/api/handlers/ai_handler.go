package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type AIHandler struct {
	aiService services.AIService
}

func NewAIHandler(aiService services.AIService) *AIHandler {
	return &AIHandler{
		aiService: aiService,
	}
}

func (h *AIHandler) AISearch(c *fiber.Ctx) error {
	var req dto.AISearchRequest
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

	result, err := h.aiService.AISearch(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "AI search failed", err)
	}

	return utils.SuccessResponse(c, "AI search completed", result)
}

func (h *AIHandler) CreateChatSession(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.CreateAIChatRequest
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

	result, err := h.aiService.CreateChatSession(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create chat session", err)
	}

	return utils.SuccessResponse(c, "Chat session created", result)
}

func (h *AIHandler) GetChatSession(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	sessionIDStr := c.Params("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid session ID")
	}

	result, err := h.aiService.GetChatSession(c.Context(), user.ID, sessionID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Chat session not found", err)
	}

	return utils.SuccessResponse(c, "Chat session retrieved", result)
}

func (h *AIHandler) GetChatSessions(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.GetAIChatSessionsRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid query parameters")
	}

	result, err := h.aiService.GetChatSessions(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get chat sessions", err)
	}

	return utils.SuccessResponse(c, "Chat sessions retrieved", result)
}

func (h *AIHandler) DeleteChatSession(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	sessionIDStr := c.Params("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid session ID")
	}

	err = h.aiService.DeleteChatSession(c.Context(), user.ID, sessionID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete chat session", err)
	}

	return utils.SuccessResponse(c, "Chat session deleted", nil)
}

func (h *AIHandler) ClearAllChatSessions(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	err = h.aiService.ClearAllChatSessions(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to clear chat sessions", err)
	}

	return utils.SuccessResponse(c, "All chat sessions cleared", nil)
}

func (h *AIHandler) SendMessage(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	sessionIDStr := c.Params("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid session ID")
	}

	var req dto.SendAIChatMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}
	req.SessionID = sessionID

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	result, err := h.aiService.SendMessage(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send message", err)
	}

	return utils.SuccessResponse(c, "Message sent", result)
}
