package routes

import (
	"github.com/gofiber/fiber/v2"

	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupAIRoutes(api fiber.Router, h *handlers.Handlers) {
	ai := api.Group("/ai")

	// Public AI search (no authentication required, but history saved if logged in)
	ai.Get("/search", middleware.OptionalAuth(), h.AIHandler.AISearch)

	// Protected AI chat endpoints
	chat := ai.Group("/chat")
	chat.Use(middleware.Protected())
	chat.Post("/", h.AIHandler.CreateChatSession)
	chat.Get("/", h.AIHandler.GetChatSessions)
	chat.Delete("/", h.AIHandler.ClearAllChatSessions)
	chat.Get("/:sessionId", h.AIHandler.GetChatSession)
	chat.Delete("/:sessionId", h.AIHandler.DeleteChatSession)
	chat.Post("/:sessionId/messages", h.AIHandler.SendMessage)
}
