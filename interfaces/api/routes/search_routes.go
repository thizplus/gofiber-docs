package routes

import (
	"github.com/gofiber/fiber/v2"

	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupSearchRoutes(api fiber.Router, h *handlers.Handlers) {
	search := api.Group("/search")

	// Public search endpoints (no authentication required, but history saved if logged in)
	search.Get("/", middleware.OptionalAuth(), h.SearchHandler.Search)
	search.Get("/websites", middleware.OptionalAuth(), h.SearchHandler.SearchWebsites)
	search.Get("/images", middleware.OptionalAuth(), h.SearchHandler.SearchImages)
	search.Get("/videos", middleware.OptionalAuth(), h.SearchHandler.SearchVideos)
	search.Get("/videos/:videoId", h.SearchHandler.GetVideoDetails)
	search.Get("/places", middleware.OptionalAuth(), h.SearchHandler.SearchPlaces)
	search.Get("/places/:placeId", h.SearchHandler.GetPlaceDetails)
	search.Get("/nearby", h.SearchHandler.SearchNearbyPlaces)

	// Protected search history endpoints
	history := search.Group("/history")
	history.Use(middleware.Protected())
	history.Get("/", h.SearchHandler.GetSearchHistory)
	history.Delete("/", h.SearchHandler.ClearSearchHistory)
	history.Delete("/:id", h.SearchHandler.DeleteSearchHistoryItem)
}
