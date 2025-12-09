# STOU Smart Tour - Backend Development Plan
# Part 5: Interface Layer (Handlers, Routes, Middleware)

---

## Table of Contents - All Parts

| Part | หัวข้อ | สถานะ |
|------|--------|-------|
| Part 1 | Project Overview & Foundation | ✅ Done |
| Part 2 | Domain Layer (Models, DTOs, Interfaces) | ✅ Done |
| Part 3 | Infrastructure Layer (External APIs, Cache) | ✅ Done |
| Part 4 | Application Layer (Services Implementation) | ✅ Done |
| **Part 5** | **Interface Layer (Handlers, Routes, Middleware)** | 📍 Current |

---

## 1. Handlers Overview

### 1.1 โครงสร้าง Handlers

```
interfaces/api/handlers/
├── handlers.go              # ✅ มีอยู่ (อัพเดท)
├── user_handler.go          # ✅ มีอยู่
├── task_handler.go          # ✅ มีอยู่
├── file_handler.go          # ✅ มีอยู่
├── job_handler.go           # ✅ มีอยู่
├── search_handler.go        # 🆕 NEW
├── ai_handler.go            # 🆕 NEW
├── folder_handler.go        # 🆕 NEW
├── favorite_handler.go      # 🆕 NEW
└── utility_handler.go       # 🆕 NEW
```

---

## 2. Search Handler

```go
// interfaces/api/handlers/search_handler.go

package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/dto"
    "github.com/your-org/stou-smart-tour/domain/services"
    "github.com/your-org/stou-smart-tour/pkg/utils"
)

type SearchHandler struct {
    searchService services.SearchService
}

func NewSearchHandler(searchService services.SearchService) *SearchHandler {
    return &SearchHandler{
        searchService: searchService,
    }
}

// Search godoc
// @Summary Search with Google Custom Search
// @Description Perform search using Google Custom Search API
// @Tags Search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param type query string false "Search type (all, website, image, video)" default(all)
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Results per page" default(10)
// @Param location query string false "Location (lat,lng)"
// @Param language query string false "Language (th, en)" default(th)
// @Success 200 {object} dto.APIResponse{data=dto.SearchResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/search [get]
func (h *SearchHandler) Search(c *fiber.Ctx) error {
    var req dto.SearchRequest
    if err := c.QueryParser(&req); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request parameters", err)
    }

    // Validate
    if err := utils.ValidateStruct(&req); err != nil {
        return utils.ValidationErrorResponse(c, err)
    }

    // Perform search
    result, err := h.searchService.Search(c.Context(), &req)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Search failed", err)
    }

    // Save search history if user is logged in
    if userID := getUserIDFromContext(c); userID != nil {
        go h.searchService.SaveSearchHistory(c.Context(), *userID, req.Query, req.Type, int(result.Meta.Total))
    }

    return utils.SuccessResponse(c, "Search completed", result)
}

// SearchPlaces godoc
// @Summary Search nearby places
// @Description Search for places near a location using Google Places API
// @Tags Search
// @Accept json
// @Produce json
// @Param lat query number true "Latitude"
// @Param lng query number true "Longitude"
// @Param radius query int false "Search radius in meters" default(5000)
// @Param type query string false "Place type (restaurant, tourist_attraction, etc.)"
// @Param keyword query string false "Keyword filter"
// @Success 200 {object} dto.APIResponse{data=dto.PlacesResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/search/places [get]
func (h *SearchHandler) SearchPlaces(c *fiber.Ctx) error {
    var req dto.PlacesSearchRequest
    if err := c.QueryParser(&req); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request parameters", err)
    }

    // Validate
    if err := utils.ValidateStruct(&req); err != nil {
        return utils.ValidationErrorResponse(c, err)
    }

    result, err := h.searchService.SearchPlaces(c.Context(), &req)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Places search failed", err)
    }

    return utils.SuccessResponse(c, "Places found", result)
}

// GetPlaceDetail godoc
// @Summary Get place details
// @Description Get detailed information about a place
// @Tags Search
// @Accept json
// @Produce json
// @Param id path string true "Place ID"
// @Success 200 {object} dto.APIResponse{data=dto.PlaceDetailResponse}
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/search/places/{id} [get]
func (h *SearchHandler) GetPlaceDetail(c *fiber.Ctx) error {
    placeID := c.Params("id")
    if placeID == "" {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Place ID is required", nil)
    }

    result, err := h.searchService.GetPlaceDetail(c.Context(), placeID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get place details", err)
    }

    return utils.SuccessResponse(c, "Place details retrieved", result)
}

// GetSearchHistory godoc
// @Summary Get search history
// @Description Get user's search history
// @Tags Search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Results per page" default(20)
// @Success 200 {object} dto.APIResponse{data=dto.SearchHistoryResponse}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/search/history [get]
func (h *SearchHandler) GetSearchHistory(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    page := c.QueryInt("page", 1)
    perPage := c.QueryInt("per_page", 20)
    offset := (page - 1) * perPage

    result, err := h.searchService.GetSearchHistory(c.Context(), *userID, offset, perPage)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get search history", err)
    }

    return utils.SuccessResponse(c, "Search history retrieved", result)
}

// ClearSearchHistory godoc
// @Summary Clear search history
// @Description Clear all search history for the current user
// @Tags Search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/search/history [delete]
func (h *SearchHandler) ClearSearchHistory(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    if err := h.searchService.ClearSearchHistory(c.Context(), *userID); err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to clear search history", err)
    }

    return utils.SuccessResponse(c, "Search history cleared", nil)
}

// Helper function to get user ID from context
func getUserIDFromContext(c *fiber.Ctx) *uuid.UUID {
    user := c.Locals("user")
    if user == nil {
        return nil
    }
    userContext, ok := user.(*dto.UserContext)
    if !ok {
        return nil
    }
    return &userContext.ID
}
```

---

## 3. AI Handler

```go
// interfaces/api/handlers/ai_handler.go

package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/dto"
    "github.com/your-org/stou-smart-tour/domain/services"
    "github.com/your-org/stou-smart-tour/pkg/utils"
)

type AIHandler struct {
    aiService services.AIService
}

func NewAIHandler(aiService services.AIService) *AIHandler {
    return &AIHandler{
        aiService: aiService,
    }
}

// AISearch godoc
// @Summary AI-enhanced search
// @Description Perform AI-enhanced search with summary generation
// @Tags AI
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param session_id query string false "Session ID for follow-up questions"
// @Success 200 {object} dto.APIResponse{data=dto.AISearchResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/search/ai [get]
func (h *AIHandler) AISearch(c *fiber.Ctx) error {
    var req dto.AISearchRequest
    if err := c.QueryParser(&req); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request parameters", err)
    }

    // Validate
    if err := utils.ValidateStruct(&req); err != nil {
        return utils.ValidationErrorResponse(c, err)
    }

    // Get user ID if logged in
    var userID *uuid.UUID
    if uid := getUserIDFromContext(c); uid != nil {
        userID = uid
    }

    result, err := h.aiService.AISearch(c.Context(), req.Query, userID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "AI search failed", err)
    }

    return utils.SuccessResponse(c, "AI search completed", result)
}

// Chat godoc
// @Summary Continue AI chat
// @Description Continue an existing AI chat conversation
// @Tags AI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AIChatRequest true "Chat request"
// @Success 200 {object} dto.APIResponse{data=dto.AIChatResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/search/ai/chat [post]
func (h *AIHandler) Chat(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    var req dto.AIChatRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
    }

    // Validate
    if err := utils.ValidateStruct(&req); err != nil {
        return utils.ValidationErrorResponse(c, err)
    }

    result, err := h.aiService.Chat(c.Context(), &req, *userID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Chat failed", err)
    }

    return utils.SuccessResponse(c, "Chat response generated", result)
}

// GetRelatedVideos godoc
// @Summary Get related videos
// @Description Get YouTube videos related to a query
// @Tags AI
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Number of videos" default(5)
// @Success 200 {object} dto.APIResponse{data=[]dto.VideoResult}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/search/ai/videos [get]
func (h *AIHandler) GetRelatedVideos(c *fiber.Ctx) error {
    query := c.Query("q")
    if query == "" {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Query is required", nil)
    }

    limit := c.QueryInt("limit", 5)
    if limit > 20 {
        limit = 20
    }

    videos, err := h.aiService.GetRelatedVideos(c.Context(), query, limit)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get videos", err)
    }

    return utils.SuccessResponse(c, "Videos retrieved", videos)
}

// GetChatSessions godoc
// @Summary Get chat sessions
// @Description Get user's AI chat sessions
// @Tags AI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Results per page" default(20)
// @Success 200 {object} dto.APIResponse{data=[]models.AIChatSession}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/search/ai/sessions [get]
func (h *AIHandler) GetChatSessions(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    page := c.QueryInt("page", 1)
    perPage := c.QueryInt("per_page", 20)
    offset := (page - 1) * perPage

    sessions, err := h.aiService.GetChatHistory(c.Context(), *userID, offset, perPage)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get sessions", err)
    }

    return utils.SuccessResponse(c, "Sessions retrieved", sessions)
}

// GetChatSession godoc
// @Summary Get chat session
// @Description Get a specific chat session with messages
// @Tags AI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Success 200 {object} dto.APIResponse{data=models.AIChatSession}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/search/ai/sessions/{id} [get]
func (h *AIHandler) GetChatSession(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    sessionID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid session ID", err)
    }

    session, err := h.aiService.GetChatSession(c.Context(), sessionID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusNotFound, "Session not found", err)
    }

    // Check ownership
    if session.UserID != *userID {
        return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied", nil)
    }

    return utils.SuccessResponse(c, "Session retrieved", session)
}

// DeleteChatSession godoc
// @Summary Delete chat session
// @Description Delete a chat session
// @Tags AI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/search/ai/sessions/{id} [delete]
func (h *AIHandler) DeleteChatSession(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    sessionID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid session ID", err)
    }

    if err := h.aiService.DeleteChatSession(c.Context(), sessionID, *userID); err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete session", err)
    }

    return utils.SuccessResponse(c, "Session deleted", nil)
}
```

---

## 4. Folder Handler

```go
// interfaces/api/handlers/folder_handler.go

package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/dto"
    "github.com/your-org/stou-smart-tour/domain/services"
    "github.com/your-org/stou-smart-tour/pkg/utils"
)

type FolderHandler struct {
    folderService services.FolderService
}

func NewFolderHandler(folderService services.FolderService) *FolderHandler {
    return &FolderHandler{
        folderService: folderService,
    }
}

// GetFolders godoc
// @Summary Get user's folders
// @Description Get all folders for the current user
// @Tags Folders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Results per page" default(20)
// @Success 200 {object} dto.APIResponse{data=dto.FolderListResponse}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/folders [get]
func (h *FolderHandler) GetFolders(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    page := c.QueryInt("page", 1)
    perPage := c.QueryInt("per_page", 20)
    offset := (page - 1) * perPage

    result, err := h.folderService.GetUserFolders(c.Context(), *userID, offset, perPage)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get folders", err)
    }

    return utils.SuccessResponse(c, "Folders retrieved", result)
}

// CreateFolder godoc
// @Summary Create a folder
// @Description Create a new folder
// @Tags Folders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateFolderRequest true "Folder data"
// @Success 201 {object} dto.APIResponse{data=dto.FolderResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/folders [post]
func (h *FolderHandler) CreateFolder(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    var req dto.CreateFolderRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
    }

    if err := utils.ValidateStruct(&req); err != nil {
        return utils.ValidationErrorResponse(c, err)
    }

    result, err := h.folderService.CreateFolder(c.Context(), &req, *userID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create folder", err)
    }

    return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
        Success: true,
        Message: "Folder created",
        Data:    result,
    })
}

// GetFolder godoc
// @Summary Get folder details
// @Description Get a folder with its items
// @Tags Folders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Folder ID"
// @Success 200 {object} dto.APIResponse{data=dto.FolderDetailResponse}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/folders/{id} [get]
func (h *FolderHandler) GetFolder(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    folderID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid folder ID", err)
    }

    result, err := h.folderService.GetFolder(c.Context(), folderID, *userID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error(), err)
    }

    return utils.SuccessResponse(c, "Folder retrieved", result)
}

// UpdateFolder godoc
// @Summary Update a folder
// @Description Update folder details
// @Tags Folders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Folder ID"
// @Param request body dto.UpdateFolderRequest true "Update data"
// @Success 200 {object} dto.APIResponse{data=dto.FolderResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/folders/{id} [put]
func (h *FolderHandler) UpdateFolder(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    folderID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid folder ID", err)
    }

    var req dto.UpdateFolderRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
    }

    result, err := h.folderService.UpdateFolder(c.Context(), folderID, &req, *userID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), err)
    }

    return utils.SuccessResponse(c, "Folder updated", result)
}

// DeleteFolder godoc
// @Summary Delete a folder
// @Description Delete a folder and all its items
// @Tags Folders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Folder ID"
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/folders/{id} [delete]
func (h *FolderHandler) DeleteFolder(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    folderID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid folder ID", err)
    }

    if err := h.folderService.DeleteFolder(c.Context(), folderID, *userID); err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), err)
    }

    return utils.SuccessResponse(c, "Folder deleted", nil)
}

// AddItem godoc
// @Summary Add item to folder
// @Description Add an item to a folder
// @Tags Folders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Folder ID"
// @Param request body dto.AddFolderItemRequest true "Item data"
// @Success 201 {object} dto.APIResponse{data=dto.FolderItemResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/folders/{id}/items [post]
func (h *FolderHandler) AddItem(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    folderID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid folder ID", err)
    }

    var req dto.AddFolderItemRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
    }

    if err := utils.ValidateStruct(&req); err != nil {
        return utils.ValidationErrorResponse(c, err)
    }

    result, err := h.folderService.AddItem(c.Context(), folderID, &req, *userID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), err)
    }

    return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
        Success: true,
        Message: "Item added to folder",
        Data:    result,
    })
}

// RemoveItem godoc
// @Summary Remove item from folder
// @Description Remove an item from a folder
// @Tags Folders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Folder ID"
// @Param itemId path string true "Item ID"
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/folders/{id}/items/{itemId} [delete]
func (h *FolderHandler) RemoveItem(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    folderID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid folder ID", err)
    }

    itemID, err := uuid.Parse(c.Params("itemId"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid item ID", err)
    }

    if err := h.folderService.RemoveItem(c.Context(), folderID, itemID, *userID); err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), err)
    }

    return utils.SuccessResponse(c, "Item removed from folder", nil)
}

// ShareFolder godoc
// @Summary Share a folder
// @Description Generate share link and QR code for a folder
// @Tags Folders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Folder ID"
// @Success 200 {object} dto.APIResponse{data=dto.FolderShareResponse}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/folders/{id}/share [post]
func (h *FolderHandler) ShareFolder(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    folderID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid folder ID", err)
    }

    result, err := h.folderService.GenerateShareLink(c.Context(), folderID, *userID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), err)
    }

    return utils.SuccessResponse(c, "Share link generated", result)
}

// GetPublicFolder godoc
// @Summary Get public folder
// @Description Get a public shared folder
// @Tags Folders
// @Accept json
// @Produce json
// @Param id path string true "Folder ID"
// @Success 200 {object} dto.APIResponse{data=dto.FolderDetailResponse}
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/shared/folder/{id} [get]
func (h *FolderHandler) GetPublicFolder(c *fiber.Ctx) error {
    folderID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid folder ID", err)
    }

    result, err := h.folderService.GetPublicFolder(c.Context(), folderID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error(), err)
    }

    return utils.SuccessResponse(c, "Folder retrieved", result)
}
```

---

## 5. Favorite Handler

```go
// interfaces/api/handlers/favorite_handler.go

package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/dto"
    "github.com/your-org/stou-smart-tour/domain/services"
    "github.com/your-org/stou-smart-tour/pkg/utils"
)

type FavoriteHandler struct {
    favoriteService services.FavoriteService
}

func NewFavoriteHandler(favoriteService services.FavoriteService) *FavoriteHandler {
    return &FavoriteHandler{
        favoriteService: favoriteService,
    }
}

// GetFavorites godoc
// @Summary Get favorites
// @Description Get user's favorites with optional type filter
// @Tags Favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string false "Filter by type (place, website, image, video)"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Results per page" default(20)
// @Success 200 {object} dto.APIResponse{data=dto.FavoriteListResponse}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/favorites [get]
func (h *FavoriteHandler) GetFavorites(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    favType := c.Query("type")
    page := c.QueryInt("page", 1)
    perPage := c.QueryInt("per_page", 20)
    offset := (page - 1) * perPage

    result, err := h.favoriteService.GetFavorites(c.Context(), *userID, favType, offset, perPage)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get favorites", err)
    }

    return utils.SuccessResponse(c, "Favorites retrieved", result)
}

// AddFavorite godoc
// @Summary Add to favorites
// @Description Add an item to favorites
// @Tags Favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AddFavoriteRequest true "Favorite data"
// @Success 201 {object} dto.APIResponse{data=dto.FavoriteResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/favorites [post]
func (h *FavoriteHandler) AddFavorite(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    var req dto.AddFavoriteRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
    }

    if err := utils.ValidateStruct(&req); err != nil {
        return utils.ValidationErrorResponse(c, err)
    }

    result, err := h.favoriteService.AddFavorite(c.Context(), &req, *userID)
    if err != nil {
        if err.Error() == "already in favorites" {
            return utils.ErrorResponse(c, fiber.StatusConflict, err.Error(), err)
        }
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to add favorite", err)
    }

    return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
        Success: true,
        Message: "Added to favorites",
        Data:    result,
    })
}

// RemoveFavorite godoc
// @Summary Remove from favorites
// @Description Remove an item from favorites
// @Tags Favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Favorite ID"
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/favorites/{id} [delete]
func (h *FavoriteHandler) RemoveFavorite(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    favoriteID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid favorite ID", err)
    }

    if err := h.favoriteService.RemoveFavorite(c.Context(), favoriteID, *userID); err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), err)
    }

    return utils.SuccessResponse(c, "Removed from favorites", nil)
}

// CheckFavorite godoc
// @Summary Check if favorited
// @Description Check if an item is in favorites
// @Tags Favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string true "Item type"
// @Param external_id query string true "External ID"
// @Success 200 {object} dto.APIResponse{data=dto.CheckFavoriteResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/favorites/check [get]
func (h *FavoriteHandler) CheckFavorite(c *fiber.Ctx) error {
    userID := getUserIDFromContext(c)
    if userID == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
    }

    favType := c.Query("type")
    externalID := c.Query("external_id")

    if favType == "" || externalID == "" {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "type and external_id are required", nil)
    }

    result, err := h.favoriteService.IsFavorited(c.Context(), *userID, favType, externalID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Check failed", err)
    }

    return utils.SuccessResponse(c, "Check completed", result)
}
```

---

## 6. Utility Handler

```go
// interfaces/api/handlers/utility_handler.go

package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/your-org/stou-smart-tour/domain/dto"
    "github.com/your-org/stou-smart-tour/domain/services"
    "github.com/your-org/stou-smart-tour/pkg/utils"
)

type UtilityHandler struct {
    translateService services.TranslateService
    qrcodeService    services.QRCodeService
}

func NewUtilityHandler(
    translateService services.TranslateService,
    qrcodeService services.QRCodeService,
) *UtilityHandler {
    return &UtilityHandler{
        translateService: translateService,
        qrcodeService:    qrcodeService,
    }
}

// Translate godoc
// @Summary Translate text
// @Description Translate text to target language
// @Tags Utilities
// @Accept json
// @Produce json
// @Param request body dto.TranslateRequest true "Translation request"
// @Success 200 {object} dto.APIResponse{data=dto.TranslateResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/translate [post]
func (h *UtilityHandler) Translate(c *fiber.Ctx) error {
    var req dto.TranslateRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
    }

    if err := utils.ValidateStruct(&req); err != nil {
        return utils.ValidationErrorResponse(c, err)
    }

    result, err := h.translateService.Translate(c.Context(), &req)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Translation failed", err)
    }

    return utils.SuccessResponse(c, "Translation completed", result)
}

// DetectLanguage godoc
// @Summary Detect language
// @Description Detect the language of text
// @Tags Utilities
// @Accept json
// @Produce json
// @Param text query string true "Text to detect"
// @Success 200 {object} dto.APIResponse{data=map[string]string}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/translate/detect [get]
func (h *UtilityHandler) DetectLanguage(c *fiber.Ctx) error {
    text := c.Query("text")
    if text == "" {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Text is required", nil)
    }

    language, err := h.translateService.DetectLanguage(c.Context(), text)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Detection failed", err)
    }

    return utils.SuccessResponse(c, "Language detected", map[string]string{
        "language": language,
    })
}

// GetSupportedLanguages godoc
// @Summary Get supported languages
// @Description Get list of supported languages for translation
// @Tags Utilities
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse{data=[]string}
// @Router /api/v1/translate/languages [get]
func (h *UtilityHandler) GetSupportedLanguages(c *fiber.Ctx) error {
    languages, err := h.translateService.GetSupportedLanguages(c.Context())
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get languages", err)
    }

    return utils.SuccessResponse(c, "Languages retrieved", languages)
}

// GenerateQRCode godoc
// @Summary Generate QR code
// @Description Generate a QR code from content
// @Tags Utilities
// @Accept json
// @Produce json
// @Param request body dto.QRCodeRequest true "QR code request"
// @Success 200 {object} dto.APIResponse{data=dto.QRCodeResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/qrcode [post]
func (h *UtilityHandler) GenerateQRCode(c *fiber.Ctx) error {
    var req dto.QRCodeRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
    }

    if err := utils.ValidateStruct(&req); err != nil {
        return utils.ValidationErrorResponse(c, err)
    }

    result, err := h.qrcodeService.Generate(c.Context(), &req)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "QR code generation failed", err)
    }

    return utils.SuccessResponse(c, "QR code generated", result)
}
```

---

## 7. Update Handlers Aggregator

```go
// interfaces/api/handlers/handlers.go

package handlers

import (
    "github.com/your-org/stou-smart-tour/domain/services"
    "github.com/your-org/stou-smart-tour/pkg/config"
)

type Handlers struct {
    // Existing handlers
    UserHandler *UserHandler
    TaskHandler *TaskHandler
    FileHandler *FileHandler
    JobHandler  *JobHandler

    // New handlers
    SearchHandler   *SearchHandler
    AIHandler       *AIHandler
    FolderHandler   *FolderHandler
    FavoriteHandler *FavoriteHandler
    UtilityHandler  *UtilityHandler
}

func NewHandlers(
    // Existing services
    userService services.UserService,
    taskService services.TaskService,
    fileService services.FileService,
    jobService services.JobService,

    // New services
    searchService services.SearchService,
    aiService services.AIService,
    folderService services.FolderService,
    favoriteService services.FavoriteService,
    translateService services.TranslateService,
    qrcodeService services.QRCodeService,

    config *config.Config,
) *Handlers {
    return &Handlers{
        // Existing
        UserHandler: NewUserHandler(userService, config),
        TaskHandler: NewTaskHandler(taskService),
        FileHandler: NewFileHandler(fileService, config),
        JobHandler:  NewJobHandler(jobService),

        // New
        SearchHandler:   NewSearchHandler(searchService),
        AIHandler:       NewAIHandler(aiService),
        FolderHandler:   NewFolderHandler(folderService),
        FavoriteHandler: NewFavoriteHandler(favoriteService),
        UtilityHandler:  NewUtilityHandler(translateService, qrcodeService),
    }
}
```

---

## 8. Routes

### 8.1 Search Routes

```go
// interfaces/api/routes/search_routes.go

package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/your-org/stou-smart-tour/interfaces/api/handlers"
    "github.com/your-org/stou-smart-tour/interfaces/api/middleware"
)

func SetupSearchRoutes(router fiber.Router, h *handlers.Handlers, m *middleware.Middleware) {
    search := router.Group("/search")

    // Public routes
    search.Get("/", m.RateLimitSearch(), h.SearchHandler.Search)
    search.Get("/places", m.RateLimitSearch(), h.SearchHandler.SearchPlaces)
    search.Get("/places/:id", h.SearchHandler.GetPlaceDetail)

    // AI routes (with higher rate limit)
    search.Get("/ai", m.RateLimitAI(), m.Optional(), h.AIHandler.AISearch)
    search.Get("/ai/videos", m.RateLimitSearch(), h.AIHandler.GetRelatedVideos)

    // Protected AI routes
    search.Post("/ai/chat", m.RateLimitAI(), m.Protected(), h.AIHandler.Chat)
    search.Get("/ai/sessions", m.Protected(), h.AIHandler.GetChatSessions)
    search.Get("/ai/sessions/:id", m.Protected(), h.AIHandler.GetChatSession)
    search.Delete("/ai/sessions/:id", m.Protected(), h.AIHandler.DeleteChatSession)

    // Protected search history
    search.Get("/history", m.Protected(), h.SearchHandler.GetSearchHistory)
    search.Delete("/history", m.Protected(), h.SearchHandler.ClearSearchHistory)
}
```

### 8.2 Folder Routes

```go
// interfaces/api/routes/folder_routes.go

package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/your-org/stou-smart-tour/interfaces/api/handlers"
    "github.com/your-org/stou-smart-tour/interfaces/api/middleware"
)

func SetupFolderRoutes(router fiber.Router, h *handlers.Handlers, m *middleware.Middleware) {
    folders := router.Group("/folders")

    // All folder routes require authentication
    folders.Use(m.Protected())

    folders.Get("/", h.FolderHandler.GetFolders)
    folders.Post("/", h.FolderHandler.CreateFolder)
    folders.Get("/:id", h.FolderHandler.GetFolder)
    folders.Put("/:id", h.FolderHandler.UpdateFolder)
    folders.Delete("/:id", h.FolderHandler.DeleteFolder)

    // Folder items
    folders.Post("/:id/items", h.FolderHandler.AddItem)
    folders.Delete("/:id/items/:itemId", h.FolderHandler.RemoveItem)

    // Sharing
    folders.Post("/:id/share", h.FolderHandler.ShareFolder)
}

func SetupPublicFolderRoutes(router fiber.Router, h *handlers.Handlers) {
    // Public shared folder access (no auth required)
    router.Get("/shared/folder/:id", h.FolderHandler.GetPublicFolder)
}
```

### 8.3 Favorite Routes

```go
// interfaces/api/routes/favorite_routes.go

package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/your-org/stou-smart-tour/interfaces/api/handlers"
    "github.com/your-org/stou-smart-tour/interfaces/api/middleware"
)

func SetupFavoriteRoutes(router fiber.Router, h *handlers.Handlers, m *middleware.Middleware) {
    favorites := router.Group("/favorites")

    // All favorite routes require authentication
    favorites.Use(m.Protected())

    favorites.Get("/", h.FavoriteHandler.GetFavorites)
    favorites.Post("/", h.FavoriteHandler.AddFavorite)
    favorites.Delete("/:id", h.FavoriteHandler.RemoveFavorite)
    favorites.Get("/check", h.FavoriteHandler.CheckFavorite)
}
```

### 8.4 Utility Routes

```go
// interfaces/api/routes/utility_routes.go

package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/your-org/stou-smart-tour/interfaces/api/handlers"
    "github.com/your-org/stou-smart-tour/interfaces/api/middleware"
)

func SetupUtilityRoutes(router fiber.Router, h *handlers.Handlers, m *middleware.Middleware) {
    // Translation routes
    translate := router.Group("/translate")
    translate.Post("/", m.RateLimitGeneral(), h.UtilityHandler.Translate)
    translate.Get("/detect", m.RateLimitGeneral(), h.UtilityHandler.DetectLanguage)
    translate.Get("/languages", h.UtilityHandler.GetSupportedLanguages)

    // QR Code routes
    qrcode := router.Group("/qrcode")
    qrcode.Post("/", m.RateLimitGeneral(), h.UtilityHandler.GenerateQRCode)
}
```

### 8.5 Update Main Routes

```go
// interfaces/api/routes/routes.go

package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/your-org/stou-smart-tour/interfaces/api/handlers"
    "github.com/your-org/stou-smart-tour/interfaces/api/middleware"
)

func SetupRoutes(app *fiber.App, h *handlers.Handlers, m *middleware.Middleware) {
    // Health check
    SetupHealthRoutes(app)

    // API v1 routes
    api := app.Group("/api/v1")

    // Existing routes
    SetupAuthRoutes(api, h, m)
    SetupUserRoutes(api, h, m)
    SetupTaskRoutes(api, h, m)
    SetupFileRoutes(api, h, m)
    SetupJobRoutes(api, h, m)

    // New routes
    SetupSearchRoutes(api, h, m)
    SetupFolderRoutes(api, h, m)
    SetupFavoriteRoutes(api, h, m)
    SetupUtilityRoutes(api, h, m)

    // Public routes (no /api/v1 prefix)
    SetupPublicFolderRoutes(app, h)

    // WebSocket (if needed)
    SetupWebSocketRoutes(app)
}
```

---

## 9. Rate Limit Middleware

```go
// interfaces/api/middleware/rate_limit_middleware.go

package middleware

import (
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    "github.com/your-org/stou-smart-tour/pkg/config"
)

type RateLimitMiddleware struct {
    config *config.Config
}

func NewRateLimitMiddleware(config *config.Config) *RateLimitMiddleware {
    return &RateLimitMiddleware{config: config}
}

// RateLimitSearch creates rate limiter for search endpoints
func (m *RateLimitMiddleware) RateLimitSearch() fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        m.config.RateLimitSearch, // default: 30 requests
        Expiration: 1 * time.Minute,
        KeyGenerator: func(c *fiber.Ctx) string {
            // Use user ID if authenticated, otherwise use IP
            if user := c.Locals("user"); user != nil {
                return "user:" + user.(*dto.UserContext).ID.String()
            }
            return "ip:" + c.IP()
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(fiber.StatusTooManyRequests).JSON(dto.ErrorResponse{
                Success: false,
                Error: dto.ErrorInfo{
                    Code:    dto.ErrCodeRateLimited,
                    Message: "Rate limit exceeded. Please try again later.",
                },
            })
        },
        SkipFailedRequests:     false,
        SkipSuccessfulRequests: false,
    })
}

// RateLimitAI creates stricter rate limiter for AI endpoints
func (m *RateLimitMiddleware) RateLimitAI() fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        m.config.RateLimitAI, // default: 10 requests
        Expiration: 1 * time.Minute,
        KeyGenerator: func(c *fiber.Ctx) string {
            if user := c.Locals("user"); user != nil {
                return "ai:user:" + user.(*dto.UserContext).ID.String()
            }
            return "ai:ip:" + c.IP()
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(fiber.StatusTooManyRequests).JSON(dto.ErrorResponse{
                Success: false,
                Error: dto.ErrorInfo{
                    Code:    dto.ErrCodeRateLimited,
                    Message: "AI rate limit exceeded. Please wait before making more AI requests.",
                },
            })
        },
    })
}

// RateLimitGeneral creates general rate limiter
func (m *RateLimitMiddleware) RateLimitGeneral() fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        m.config.RateLimitGeneral, // default: 100 requests
        Expiration: 1 * time.Minute,
        KeyGenerator: func(c *fiber.Ctx) string {
            if user := c.Locals("user"); user != nil {
                return "general:user:" + user.(*dto.UserContext).ID.String()
            }
            return "general:ip:" + c.IP()
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(fiber.StatusTooManyRequests).JSON(dto.ErrorResponse{
                Success: false,
                Error: dto.ErrorInfo{
                    Code:    dto.ErrCodeRateLimited,
                    Message: "Too many requests. Please slow down.",
                },
            })
        },
    })
}
```

### 9.1 Update Middleware Aggregator

```go
// interfaces/api/middleware/middleware.go

package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/your-org/stou-smart-tour/pkg/config"
)

type Middleware struct {
    auth      *AuthMiddleware
    rateLimit *RateLimitMiddleware
    config    *config.Config
}

func NewMiddleware(config *config.Config) *Middleware {
    return &Middleware{
        auth:      NewAuthMiddleware(config),
        rateLimit: NewRateLimitMiddleware(config),
        config:    config,
    }
}

// Auth middleware methods
func (m *Middleware) Protected() fiber.Handler {
    return m.auth.Protected()
}

func (m *Middleware) Optional() fiber.Handler {
    return m.auth.Optional()
}

func (m *Middleware) AdminOnly() fiber.Handler {
    return m.auth.AdminOnly()
}

func (m *Middleware) RequireRole(role string) fiber.Handler {
    return m.auth.RequireRole(role)
}

// Rate limit middleware methods
func (m *Middleware) RateLimitSearch() fiber.Handler {
    return m.rateLimit.RateLimitSearch()
}

func (m *Middleware) RateLimitAI() fiber.Handler {
    return m.rateLimit.RateLimitAI()
}

func (m *Middleware) RateLimitGeneral() fiber.Handler {
    return m.rateLimit.RateLimitGeneral()
}
```

---

## 10. Update Main Entry Point

```go
// cmd/api/main.go

package main

import (
    "log"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"

    "github.com/your-org/stou-smart-tour/interfaces/api/handlers"
    "github.com/your-org/stou-smart-tour/interfaces/api/middleware"
    "github.com/your-org/stou-smart-tour/interfaces/api/routes"
    "github.com/your-org/stou-smart-tour/pkg/di"
)

func main() {
    // Initialize DI container
    container := di.NewContainer()

    // Create Fiber app
    app := fiber.New(fiber.Config{
        ErrorHandler: middleware.ErrorHandler(),
    })

    // Global middleware
    app.Use(recover.New())
    app.Use(logger.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins:     "*",
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
        AllowCredentials: true,
    }))

    // Initialize handlers
    h := handlers.NewHandlers(
        // Existing services
        container.UserService,
        container.TaskService,
        container.FileService,
        container.JobService,

        // New services
        container.SearchService,
        container.AIService,
        container.FolderService,
        container.FavoriteService,
        container.TranslateService,
        container.QRCodeService,

        container.Config,
    )

    // Initialize middleware
    m := middleware.NewMiddleware(container.Config)

    // Setup routes
    routes.SetupRoutes(app, h, m)

    // Start server
    port := container.Config.AppPort
    log.Printf("Server starting on port %s", port)
    if err := app.Listen(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

---

## 11. Summary - Interface Layer Files

```
interfaces/api/
├── handlers/
│   ├── handlers.go              # ✏️ UPDATE (add new handlers)
│   ├── user_handler.go          # ✅ KEEP
│   ├── task_handler.go          # ✅ KEEP
│   ├── file_handler.go          # ✅ KEEP
│   ├── job_handler.go           # ✅ KEEP
│   ├── search_handler.go        # 🆕 NEW
│   ├── ai_handler.go            # 🆕 NEW
│   ├── folder_handler.go        # 🆕 NEW
│   ├── favorite_handler.go      # 🆕 NEW
│   └── utility_handler.go       # 🆕 NEW
│
├── middleware/
│   ├── auth_middleware.go       # ✅ KEEP
│   ├── cors_middleware.go       # ✅ KEEP
│   ├── error_middleware.go      # ✅ KEEP
│   ├── logger_middleware.go     # ✅ KEEP
│   ├── rate_limit_middleware.go # 🆕 NEW
│   └── middleware.go            # 🆕 NEW (aggregator)
│
└── routes/
    ├── routes.go                # ✏️ UPDATE
    ├── auth_routes.go           # ✅ KEEP
    ├── user_routes.go           # ✅ KEEP
    ├── task_routes.go           # ✅ KEEP
    ├── file_routes.go           # ✅ KEEP
    ├── job_routes.go            # ✅ KEEP
    ├── health_routes.go         # ✅ KEEP
    ├── websocket_routes.go      # ✅ KEEP
    ├── search_routes.go         # 🆕 NEW
    ├── folder_routes.go         # 🆕 NEW
    ├── favorite_routes.go       # 🆕 NEW
    └── utility_routes.go        # 🆕 NEW
```

---

## 12. API Endpoints Summary

### 12.1 Complete API Reference

```
┌─────────────────────────────────────────────────────────────────┐
│                    STOU Smart Tour API                           │
│                    Complete Endpoints                            │
└─────────────────────────────────────────────────────────────────┘

Authentication (Existing)
├── POST   /api/v1/auth/register          # Register new user
├── POST   /api/v1/auth/login             # User login
└── GET    /api/v1/auth/me                # Get current user

Users (Existing)
├── GET    /api/v1/users/profile          # Get profile
├── PUT    /api/v1/users/profile          # Update profile
└── GET    /api/v1/users/                 # List users (admin)

Search (NEW)
├── GET    /api/v1/search                 # Google search
│          ?q={query}&type={type}&page={page}
├── GET    /api/v1/search/places          # Nearby places
│          ?lat={lat}&lng={lng}&radius={radius}
├── GET    /api/v1/search/places/:id      # Place details
├── GET    /api/v1/search/history         # Search history (auth)
└── DELETE /api/v1/search/history         # Clear history (auth)

AI Mode (NEW)
├── GET    /api/v1/search/ai              # AI search + summary
│          ?q={query}
├── POST   /api/v1/search/ai/chat         # Continue chat (auth)
├── GET    /api/v1/search/ai/videos       # Related videos
│          ?q={query}&limit={limit}
├── GET    /api/v1/search/ai/sessions     # Chat sessions (auth)
├── GET    /api/v1/search/ai/sessions/:id # Get session (auth)
└── DELETE /api/v1/search/ai/sessions/:id # Delete session (auth)

Folders (NEW) - All require authentication
├── GET    /api/v1/folders                # List folders
├── POST   /api/v1/folders                # Create folder
├── GET    /api/v1/folders/:id            # Get folder + items
├── PUT    /api/v1/folders/:id            # Update folder
├── DELETE /api/v1/folders/:id            # Delete folder
├── POST   /api/v1/folders/:id/items      # Add item
├── DELETE /api/v1/folders/:id/items/:itemId  # Remove item
└── POST   /api/v1/folders/:id/share      # Generate share link

Public Folder (NEW)
└── GET    /shared/folder/:id             # View public folder

Favorites (NEW) - All require authentication
├── GET    /api/v1/favorites              # List favorites
│          ?type={type}
├── POST   /api/v1/favorites              # Add favorite
├── DELETE /api/v1/favorites/:id          # Remove favorite
└── GET    /api/v1/favorites/check        # Check if favorited
           ?type={type}&external_id={id}

Translation (NEW)
├── POST   /api/v1/translate              # Translate text
├── GET    /api/v1/translate/detect       # Detect language
│          ?text={text}
└── GET    /api/v1/translate/languages    # Supported languages

QR Code (NEW)
└── POST   /api/v1/qrcode                 # Generate QR code

Health (Existing)
├── GET    /health                        # Health check
└── GET    /                              # Welcome
```

---

## 13. Final Project Structure

```
gofiber-docs/
├── cmd/
│   └── api/
│       └── main.go                       # ✏️ UPDATE
│
├── domain/
│   ├── models/
│   │   ├── user.go                       # ✏️ UPDATE
│   │   ├── folder.go                     # 🆕 NEW
│   │   ├── folder_item.go                # 🆕 NEW
│   │   ├── favorite.go                   # 🆕 NEW
│   │   ├── search_history.go             # 🆕 NEW
│   │   ├── ai_chat_session.go            # 🆕 NEW
│   │   └── ai_chat_message.go            # 🆕 NEW
│   ├── repositories/                     # + 5 new interfaces
│   ├── services/                         # + 6 new interfaces
│   └── dto/                              # + 5 new files
│
├── application/
│   └── serviceimpl/                      # + 6 new implementations
│
├── infrastructure/
│   ├── postgres/                         # + 5 new repo impls
│   ├── redis/                            # ✏️ UPDATE
│   ├── cache/                            # 🆕 NEW folder
│   └── external/                         # 🆕 NEW folder
│       ├── google/                       # 5 clients
│       └── openai/                       # 1 client
│
├── interfaces/api/
│   ├── handlers/                         # + 5 new handlers
│   ├── middleware/                       # + 2 new middlewares
│   └── routes/                           # + 4 new route files
│
└── pkg/
    ├── config/                           # ✏️ UPDATE
    └── di/                               # ✏️ UPDATE

Total New Files: ~40 files
Total Updated Files: ~10 files
```

---

## 14. Development Checklist

```
Phase 1: Foundation
□ Update User model (add student_id)
□ Add new config fields
□ Setup external API clients (Google, OpenAI)
□ Update DI container
□ Create cache keys

Phase 2: Core Features
□ Search service + handler + routes
□ Folder model + repository + service + handler
□ Folder item model + repository
□ Favorite model + repository + service + handler

Phase 3: Advanced Features
□ AI service + handler
□ OpenAI integration
□ YouTube integration
□ AI chat sessions

Phase 4: Utilities & Polish
□ Translate service + handler
□ QR code service + handler
□ Rate limit middleware
□ Testing
□ Documentation
```

---

## Congratulations! 🎉

เอกสารแผนการพัฒนา STOU Smart Tour Backend ครบทั้ง 5 Parts แล้ว!

### Summary of All Parts:

| Part | File Name | Content |
|------|-----------|---------|
| 1 | `STOU-Backend-Plan-Part1-Overview.md` | Project Overview, Tech Stack, Structure |
| 2 | `STOU-Backend-Plan-Part2-Domain.md` | Models, DTOs, Repository & Service Interfaces |
| 3 | `STOU-Backend-Plan-Part3-Infrastructure.md` | External APIs, Cache, Repository Implementations |
| 4 | `STOU-Backend-Plan-Part4-Application.md` | Service Implementations, DI Container |
| 5 | `STOU-Backend-Plan-Part5-Interface.md` | Handlers, Routes, Middleware |

---

*Document Version: 1.0*
*Part: 5 of 5*
*Status: COMPLETE*
