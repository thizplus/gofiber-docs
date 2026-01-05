package routes

import (
	"github.com/gofiber/fiber/v2"

	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupAIRoutes(api fiber.Router, h *handlers.Handlers) {
	ai := api.Group("/ai")

	// All AI features require login
	ai.Use(middleware.Protected())

	// AI search with summary
	ai.Get("/search", h.AIHandler.AISearch)

	// AI chat endpoints
	chat := ai.Group("/chat")
	chat.Post("/", h.AIHandler.CreateChatSession)
	chat.Get("/", h.AIHandler.GetChatSessions)
	chat.Delete("/", h.AIHandler.ClearAllChatSessions)
	chat.Get("/:sessionId", h.AIHandler.GetChatSession)
	chat.Delete("/:sessionId", h.AIHandler.DeleteChatSession)
	chat.Post("/:sessionId/messages", h.AIHandler.SendMessage)
}
