package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
)

func SetupAuthRoutes(api fiber.Router, h *handlers.Handlers) {
	auth := api.Group("/auth")

	// Traditional auth (email/password)
	auth.Post("/register", h.UserHandler.Register)
	auth.Post("/login", h.UserHandler.Login)

	// Google OAuth
	auth.Get("/google", h.UserHandler.GoogleAuth)
	auth.Get("/google/callback", h.UserHandler.GoogleCallback)

	// LINE OAuth
	auth.Get("/line", h.UserHandler.LineAuth)
	auth.Get("/line/callback", h.UserHandler.LineCallback)
}