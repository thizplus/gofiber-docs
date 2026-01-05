package routes

import (
	"github.com/gofiber/fiber/v2"

	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupSearchRoutes(api fiber.Router, h *handlers.Handlers) {
	search := api.Group("/search")

	// Initialize guest rate limiter
	guestLimiter := middleware.NewGuestRateLimitMiddleware(middleware.DefaultGuestRateLimitConfig())

	// Public search endpoints with guest rate limiting
	// OptionalAuth allows logged users to save search history
	search.Get("/", middleware.OptionalAuth(), guestLimiter.GuestSearchLimit(), h.SearchHandler.Search)
	search.Get("/websites", middleware.OptionalAuth(), guestLimiter.GuestSearchLimit(), h.SearchHandler.SearchWebsites)
	search.Get("/images", middleware.OptionalAuth(), guestLimiter.GuestMediaLimit(), h.SearchHandler.SearchImages)
	search.Get("/videos", middleware.OptionalAuth(), guestLimiter.GuestMediaLimit(), h.SearchHandler.SearchVideos)
	search.Get("/videos/:videoId", h.SearchHandler.GetVideoDetails)
	search.Get("/places", middleware.OptionalAuth(), h.SearchHandler.SearchPlaces) // No rate limit - public
	search.Get("/places/:placeId", h.SearchHandler.GetPlaceDetails)
	search.Get("/places/:placeId/enhanced", middleware.OptionalAuth(), h.SearchHandler.GetPlaceDetailsEnhanced)
	search.Get("/nearby", middleware.OptionalAuth(), guestLimiter.GuestPlacesLimit(), h.SearchHandler.SearchNearbyPlaces)

	// Protected search history endpoints (login required)
	history := search.Group("/history")
	history.Use(middleware.Protected())
	history.Get("/", h.SearchHandler.GetSearchHistory)
	history.Delete("/", h.SearchHandler.ClearSearchHistory)
	history.Delete("/:id", h.SearchHandler.DeleteSearchHistoryItem)
}
