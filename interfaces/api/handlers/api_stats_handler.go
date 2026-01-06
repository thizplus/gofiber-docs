package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"gofiber-template/application/serviceimpl"
	"gofiber-template/pkg/utils"
)

type APIStatsHandler struct {
	logger *serviceimpl.APILoggerService
}

func NewAPIStatsHandler(logger *serviceimpl.APILoggerService) *APIStatsHandler {
	return &APIStatsHandler{logger: logger}
}

// GetSummary returns API usage summary
// @Summary Get API usage summary
// @Tags Admin
// @Accept json
// @Produce json
// @Param days query int false "Number of days (default: 7)"
// @Success 200 {object} map[string]interface{}
// @Router /admin/api-stats/summary [get]
func (h *APIStatsHandler) GetSummary(c *fiber.Ctx) error {
	days := c.QueryInt("days", 7)
	if days < 1 {
		days = 7
	}
	if days > 90 {
		days = 90
	}

	summary, err := h.logger.GetSummary(c.Context(), days)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get summary", err)
	}

	return utils.SuccessResponse(c, "API usage summary", summary)
}

// GetServiceStats returns statistics grouped by service
// @Summary Get API statistics by service
// @Tags Admin
// @Accept json
// @Produce json
// @Param days query int false "Number of days (default: 7)"
// @Success 200 {object} map[string]interface{}
// @Router /admin/api-stats/services [get]
func (h *APIStatsHandler) GetServiceStats(c *fiber.Ctx) error {
	days := c.QueryInt("days", 7)
	if days < 1 {
		days = 7
	}

	stats, err := h.logger.GetStatsByService(c.Context(), days)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get stats", err)
	}

	return utils.SuccessResponse(c, "Service statistics", fiber.Map{
		"period": strconv.Itoa(days) + " days",
		"stats":  stats,
	})
}

// GetEndpointStats returns statistics grouped by endpoint
// @Summary Get API statistics by endpoint
// @Tags Admin
// @Accept json
// @Produce json
// @Param service query string false "Service name filter"
// @Param days query int false "Number of days (default: 7)"
// @Success 200 {object} map[string]interface{}
// @Router /admin/api-stats/endpoints [get]
func (h *APIStatsHandler) GetEndpointStats(c *fiber.Ctx) error {
	service := c.Query("service", "")
	days := c.QueryInt("days", 7)
	if days < 1 {
		days = 7
	}

	stats, err := h.logger.GetStatsByEndpoint(c.Context(), service, days)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get stats", err)
	}

	return utils.SuccessResponse(c, "Endpoint statistics", fiber.Map{
		"period":  strconv.Itoa(days) + " days",
		"service": service,
		"stats":   stats,
	})
}

// GetDailyStats returns daily statistics
// @Summary Get daily API statistics
// @Tags Admin
// @Accept json
// @Produce json
// @Param days query int false "Number of days (default: 30)"
// @Success 200 {object} map[string]interface{}
// @Router /admin/api-stats/daily [get]
func (h *APIStatsHandler) GetDailyStats(c *fiber.Ctx) error {
	days := c.QueryInt("days", 30)
	if days < 1 {
		days = 30
	}
	if days > 90 {
		days = 90
	}

	stats, err := h.logger.GetDailyStats(c.Context(), days)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get stats", err)
	}

	return utils.SuccessResponse(c, "Daily statistics", fiber.Map{
		"period": strconv.Itoa(days) + " days",
		"stats":  stats,
	})
}

// GetCostBreakdown returns cost breakdown by service
// @Summary Get API cost breakdown
// @Tags Admin
// @Accept json
// @Produce json
// @Param days query int false "Number of days (default: 30)"
// @Success 200 {object} map[string]interface{}
// @Router /admin/api-stats/costs [get]
func (h *APIStatsHandler) GetCostBreakdown(c *fiber.Ctx) error {
	days := c.QueryInt("days", 30)
	if days < 1 {
		days = 30
	}

	costs, err := h.logger.GetServiceCosts(c.Context(), days)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get costs", err)
	}

	totalCost, _ := h.logger.GetTotalCost(c.Context(), days)
	cacheHitRate, _ := h.logger.GetCacheHitRate(c.Context(), days)

	return utils.SuccessResponse(c, "Cost breakdown", fiber.Map{
		"period":       strconv.Itoa(days) + " days",
		"totalCost":    totalCost,
		"cacheHitRate": cacheHitRate,
		"costSaved":    totalCost * (cacheHitRate / 100),
		"breakdown":    costs,
	})
}

// CleanupOldLogs removes old log entries
// @Summary Cleanup old API logs
// @Tags Admin
// @Accept json
// @Produce json
// @Param days query int false "Keep logs for last N days (default: 30)"
// @Success 200 {object} map[string]interface{}
// @Router /admin/api-stats/cleanup [delete]
func (h *APIStatsHandler) CleanupOldLogs(c *fiber.Ctx) error {
	days := c.QueryInt("days", 30)
	if days < 7 {
		days = 7 // Minimum 7 days retention
	}

	deleted, err := h.logger.CleanupOldLogs(c.Context(), days)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to cleanup", err)
	}

	return utils.SuccessResponse(c, "Cleanup completed", fiber.Map{
		"deletedCount":  deleted,
		"retentionDays": days,
	})
}
