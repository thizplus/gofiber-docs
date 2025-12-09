package routes

import (
	"github.com/gofiber/fiber/v2"

	"gofiber-template/interfaces/api/handlers"
)

func SetupUtilityRoutes(api fiber.Router, h *handlers.Handlers) {
	utils := api.Group("/utils")

	// Public config for frontend
	utils.Get("/config", h.UtilityHandler.GetPublicConfig)

	// Translation endpoints (public)
	utils.Post("/translate", h.UtilityHandler.Translate)
	utils.Post("/detect-language", h.UtilityHandler.DetectLanguage)

	// QR Code generation (public)
	utils.Post("/qrcode", h.UtilityHandler.GenerateQRCode)

	// Distance calculation (public)
	utils.Get("/distance", h.UtilityHandler.CalculateDistance)
}
