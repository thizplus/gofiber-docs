package routes

import (
	"github.com/gofiber/fiber/v2"

	"gofiber-template/application/serviceimpl"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

// SetupAdminRoutes sets up admin routes for API statistics
func SetupAdminRoutes(api fiber.Router, logger *serviceimpl.APILoggerService) {
	statsHandler := handlers.NewAPIStatsHandler(logger)

	// Admin routes - require authentication and admin role
	admin := api.Group("/admin")
	admin.Use(middleware.Protected())
	admin.Use(middleware.AdminOnly())

	// API Statistics routes
	stats := admin.Group("/api-stats")
	stats.Get("/summary", statsHandler.GetSummary)
	stats.Get("/services", statsHandler.GetServiceStats)
	stats.Get("/endpoints", statsHandler.GetEndpointStats)
	stats.Get("/daily", statsHandler.GetDailyStats)
	stats.Get("/costs", statsHandler.GetCostBreakdown)
	stats.Delete("/cleanup", statsHandler.CleanupOldLogs)
}
