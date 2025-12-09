package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
)

func SetupAuthRoutes(api fiber.Router, h *handlers.Handlers) {
	auth := api.Group("/auth")
	auth.Post("/register", h.UserHandler.Register)
	auth.Post("/login", h.UserHandler.Login)
}