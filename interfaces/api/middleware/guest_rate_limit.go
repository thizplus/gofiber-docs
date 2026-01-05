package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type GuestRateLimitConfig struct {
	SearchLimit int // per day
	PlacesLimit int // per day
	MediaLimit  int // per day
}

type GuestRateLimitMiddleware struct {
	config GuestRateLimitConfig
}

func NewGuestRateLimitMiddleware(cfg GuestRateLimitConfig) *GuestRateLimitMiddleware {
	return &GuestRateLimitMiddleware{config: cfg}
}

// DefaultGuestRateLimitConfig returns default rate limit config
func DefaultGuestRateLimitConfig() GuestRateLimitConfig {
	return GuestRateLimitConfig{
		SearchLimit: 10, // 10 searches per day
		PlacesLimit: 5,  // 5 place searches per day
		MediaLimit:  5,  // 5 media searches per day
	}
}

// GuestSearchLimit limits search for non-logged users
func (m *GuestRateLimitMiddleware) GuestSearchLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        m.config.SearchLimit,
		Expiration: 24 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			// If logged in, use different key to not apply limit
			if c.Locals("user") != nil {
				return "user:" + c.IP() // Different key for logged users
			}
			return "guest:search:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "RATE_LIMITED",
					"message": "คุณค้นหาเกินจำนวนที่กำหนดสำหรับวันนี้ กรุณาเข้าสู่ระบบเพื่อค้นหาเพิ่มเติม",
					"action":  "login_required",
				},
			})
		},
		SkipFailedRequests:     true,
		SkipSuccessfulRequests: false,
		// Skip rate limiting if user is logged in
		Next: func(c *fiber.Ctx) bool {
			return c.Locals("user") != nil
		},
	})
}

// GuestPlacesLimit limits places search for guests
func (m *GuestRateLimitMiddleware) GuestPlacesLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        m.config.PlacesLimit,
		Expiration: 24 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			if c.Locals("user") != nil {
				return "user:places:" + c.IP()
			}
			return "guest:places:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "RATE_LIMITED",
					"message": "คุณค้นหาสถานที่เกินจำนวนที่กำหนดสำหรับวันนี้ กรุณาเข้าสู่ระบบเพื่อค้นหาเพิ่มเติม",
					"action":  "login_required",
				},
			})
		},
		Next: func(c *fiber.Ctx) bool {
			return c.Locals("user") != nil
		},
	})
}

// GuestMediaLimit limits media search for guests
func (m *GuestRateLimitMiddleware) GuestMediaLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        m.config.MediaLimit,
		Expiration: 24 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			if c.Locals("user") != nil {
				return "user:media:" + c.IP()
			}
			return "guest:media:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "RATE_LIMITED",
					"message": "คุณค้นหารูปภาพ/วิดีโอเกินจำนวนที่กำหนดสำหรับวันนี้ กรุณาเข้าสู่ระบบเพื่อค้นหาเพิ่มเติม",
					"action":  "login_required",
				},
			})
		},
		Next: func(c *fiber.Ctx) bool {
			return c.Locals("user") != nil
		},
	})
}
