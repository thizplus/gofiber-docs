package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
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

	user, err := h.userService.Register(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Registration failed", err)
	}

	// Generate token for auto-login after registration
	token, err := h.userService.GenerateJWT(user)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token", err)
	}

	// Return same format as login for auto-login
	registerResponse := &dto.LoginResponse{
		Token: token,
		User:  *dto.UserToUserResponse(user),
	}
	return utils.SuccessResponse(c, "ลงทะเบียนสำเร็จ", registerResponse)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
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

	token, user, err := h.userService.Login(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Login failed", err)
	}

	loginResponse := &dto.LoginResponse{
		Token: token,
		User:  *dto.UserToUserResponse(user),
	}
	return utils.SuccessResponse(c, "Login successful", loginResponse)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	profile, err := h.userService.GetProfile(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found", err)
	}

	profileResponse := dto.UserToUserResponse(profile)
	return utils.SuccessResponse(c, "Profile retrieved successfully", profileResponse)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	updatedUser, err := h.userService.UpdateProfile(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Profile update failed", err)
	}

	userResponse := dto.UserToUserResponse(updatedUser)
	return utils.SuccessResponse(c, "Profile updated successfully", userResponse)
}

// UpdateProfileInfo updates user profile information (firstName, lastName, studentId, language, theme)
func (h *UserHandler) UpdateProfileInfo(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.UpdateProfileRequest
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

	updatedUser, err := h.userService.UpdateProfileInfo(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	userResponse := dto.UserToUserResponse(updatedUser)
	return utils.SuccessResponse(c, "อัปเดตข้อมูลสำเร็จ", userResponse)
}

// UpdateAvatar uploads a new avatar image
func (h *UserHandler) UpdateAvatar(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Get file from form
	file, err := c.FormFile("avatar")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "กรุณาเลือกไฟล์รูปภาพ", err)
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ไม่สามารถอ่านไฟล์ได้", err)
	}
	defer src.Close()

	// Read file content
	fileData := make([]byte, file.Size)
	if _, err := src.Read(fileData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ไม่สามารถอ่านไฟล์ได้", err)
	}

	// Get content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	// Update avatar
	result, err := h.userService.UpdateAvatar(c.Context(), user.ID, fileData, contentType)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	return utils.SuccessResponse(c, "อัปโหลดรูปโปรไฟล์สำเร็จ", result)
}

// DeleteAvatar removes user avatar
func (h *UserHandler) DeleteAvatar(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	if err := h.userService.DeleteAvatar(c.Context(), user.ID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), err)
	}

	return utils.SuccessResponse(c, "ลบรูปโปรไฟล์สำเร็จ", nil)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	err = h.userService.DeleteUser(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User deletion failed", err)
	}

	return utils.SuccessResponse(c, "User deleted successfully", nil)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	offsetStr := c.Query("offset", "0")
	limitStr := c.Query("limit", "10")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid offset parameter")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid limit parameter")
	}

	users, total, err := h.userService.ListUsers(c.Context(), offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve users", err)
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *dto.UserToUserResponse(user)
	}

	response := &dto.UserListResponse{
		Users: userResponses,
		Meta: dto.PaginationMeta{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}

	return utils.SuccessResponse(c, "Users retrieved successfully", response)
}

// ==================== OAuth Handlers ====================

// GoogleAuth redirects to Google OAuth
func (h *UserHandler) GoogleAuth(c *fiber.Ctx) error {
	frontendURL := c.Query("frontend_url", os.Getenv("FRONTEND_URL"))

	// Create state data
	stateData := dto.OAuthStateData{
		FrontendURL: frontendURL,
	}

	// Encode state as base64
	stateJSON, err := json.Marshal(stateData)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create OAuth state", err)
	}
	state := base64.URLEncoding.EncodeToString(stateJSON)

	// Get authorization URL
	authURL := h.userService.GetGoogleAuthURL(state)

	return c.Redirect(authURL, fiber.StatusTemporaryRedirect)
}

// GoogleCallback handles Google OAuth callback
func (h *UserHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Authorization code is required", nil)
	}

	// Decode state
	stateJSON, err := base64.URLEncoding.DecodeString(state)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid OAuth state", err)
	}

	var stateData dto.OAuthStateData
	if err := json.Unmarshal(stateJSON, &stateData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid OAuth state data", err)
	}

	// Get frontend URL
	frontendURL := stateData.FrontendURL
	if frontendURL == "" {
		frontendURL = os.Getenv("FRONTEND_URL")
	}

	// Handle callback
	token, user, isNewUser, err := h.userService.HandleGoogleCallback(c.Context(), code, stateData)
	if err != nil {
		// Redirect to frontend with error
		errorURL := fmt.Sprintf("%s/auth/callback?error=%s", frontendURL, url.QueryEscape(err.Error()))
		return c.Redirect(errorURL, fiber.StatusTemporaryRedirect)
	}

	// Create callback URL with token
	callbackURL := fmt.Sprintf("%s/auth/callback?token=%s&is_new_user=%t&user_id=%s",
		frontendURL,
		url.QueryEscape(token),
		isNewUser,
		user.ID.String(),
	)

	return c.Redirect(callbackURL, fiber.StatusTemporaryRedirect)
}

// LineAuth redirects to LINE OAuth
func (h *UserHandler) LineAuth(c *fiber.Ctx) error {
	frontendURL := c.Query("frontend_url", os.Getenv("FRONTEND_URL"))

	// Create state data
	stateData := dto.OAuthStateData{
		FrontendURL: frontendURL,
	}

	// Encode state as base64
	stateJSON, err := json.Marshal(stateData)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create OAuth state", err)
	}
	state := base64.URLEncoding.EncodeToString(stateJSON)

	// Get authorization URL
	authURL := h.userService.GetLineAuthURL(state)

	return c.Redirect(authURL, fiber.StatusTemporaryRedirect)
}

// LineCallback handles LINE OAuth callback
func (h *UserHandler) LineCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Authorization code is required", nil)
	}

	// Decode state
	stateJSON, err := base64.URLEncoding.DecodeString(state)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid OAuth state", err)
	}

	var stateData dto.OAuthStateData
	if err := json.Unmarshal(stateJSON, &stateData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid OAuth state data", err)
	}

	// Get frontend URL
	frontendURL := stateData.FrontendURL
	if frontendURL == "" {
		frontendURL = os.Getenv("FRONTEND_URL")
	}

	// Handle callback
	token, user, isNewUser, err := h.userService.HandleLineCallback(c.Context(), code, stateData)
	if err != nil {
		// Redirect to frontend with error
		errorURL := fmt.Sprintf("%s/auth/callback?error=%s", frontendURL, url.QueryEscape(err.Error()))
		return c.Redirect(errorURL, fiber.StatusTemporaryRedirect)
	}

	// Create callback URL with token
	callbackURL := fmt.Sprintf("%s/auth/callback?token=%s&is_new_user=%t&user_id=%s",
		frontendURL,
		url.QueryEscape(token),
		isNewUser,
		user.ID.String(),
	)

	return c.Redirect(callbackURL, fiber.StatusTemporaryRedirect)
}