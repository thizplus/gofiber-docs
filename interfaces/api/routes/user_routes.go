package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupUserRoutes(api fiber.Router, h *handlers.Handlers) {
	users := api.Group("/users")
	users.Use(middleware.Protected())

	// Profile
	users.Get("/profile", h.UserHandler.GetProfile)
	users.Put("/profile", h.UserHandler.UpdateProfile)
	users.Patch("/profile", h.UserHandler.UpdateProfileInfo) // For partial updates (firstName, lastName, studentId, etc.)
	users.Delete("/profile", h.UserHandler.DeleteUser)

	// Avatar
	users.Post("/avatar", h.UserHandler.UpdateAvatar)
	users.Delete("/avatar", h.UserHandler.DeleteAvatar)

	// Admin only
	users.Get("/", middleware.AdminOnly(), h.UserHandler.ListUsers)
}