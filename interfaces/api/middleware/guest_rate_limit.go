package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"gofiber-template/domain/models"
)

type GuestRateLimitConfig struct {
	SearchLimit int // per day for guests
	PlacesLimit int // per day for guests
	MediaLimit  int // per day for guests
}

type AuthRateLimitConfig struct {
	SearchLimit int // per day for authenticated users
	PlacesLimit int // per day for authenticated users
	MediaLimit  int // per day for authenticated users
}

type RateLimitMiddleware struct {
	guestConfig GuestRateLimitConfig
	authConfig  AuthRateLimitConfig
}

func NewRateLimitMiddleware(guestCfg GuestRateLimitConfig, authCfg AuthRateLimitConfig) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		guestConfig: guestCfg,
		authConfig:  authCfg,
	}
}

// DefaultGuestRateLimitConfig returns default rate limit config for guests
func DefaultGuestRateLimitConfig() GuestRateLimitConfig {
	return GuestRateLimitConfig{
		SearchLimit: 10, // 10 searches per day
		PlacesLimit: 5,  // 5 place searches per day
		MediaLimit:  5,  // 5 media searches per day
	}
}

// DefaultAuthRateLimitConfig returns default rate limit config for authenticated users
func DefaultAuthRateLimitConfig() AuthRateLimitConfig {
	return AuthRateLimitConfig{
		SearchLimit: 200, // 200 searches per day
		PlacesLimit: 100, // 100 place detail views per day
		MediaLimit:  100, // 100 media searches per day
	}
}

// getUserID gets user ID from context if authenticated
func getUserID(c *fiber.Ctx) string {
	user := c.Locals("user")
	if user == nil {
		return ""
	}
	if u, ok := user.(*models.User); ok {
		return u.ID.String()
	}
	return ""
}

// SearchLimit limits search for both guests and authenticated users
func (m *RateLimitMiddleware) SearchLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max: m.authConfig.SearchLimit, // Use auth limit as max
		Expiration: 24 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			userID := getUserID(c)
			if userID != "" {
				return "auth:search:" + userID
			}
			return "guest:search:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			userID := getUserID(c)
			if userID != "" {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"success": false,
					"error": fiber.Map{
						"code":    "RATE_LIMITED",
						"message": "คุณค้นหาเกินจำนวนที่กำหนดสำหรับวันนี้ (200 ครั้ง/วัน) กรุณาลองใหม่พรุ่งนี้",
						"action":  "wait",
					},
				})
			}
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
		LimiterMiddleware: limiter.FixedWindow{},
		// Dynamic max based on user type
		Next: func(c *fiber.Ctx) bool {
			return false // Never skip - apply to all
		},
	})
}

// PlacesLimit limits places search for both guests and authenticated users
func (m *RateLimitMiddleware) PlacesLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max: m.authConfig.PlacesLimit,
		Expiration: 24 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			userID := getUserID(c)
			if userID != "" {
				return "auth:places:" + userID
			}
			return "guest:places:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			userID := getUserID(c)
			if userID != "" {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"success": false,
					"error": fiber.Map{
						"code":    "RATE_LIMITED",
						"message": "คุณดูรายละเอียดสถานที่เกินจำนวนที่กำหนดสำหรับวันนี้ (100 ครั้ง/วัน)",
						"action":  "wait",
					},
				})
			}
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "RATE_LIMITED",
					"message": "คุณค้นหาสถานที่เกินจำนวนที่กำหนดสำหรับวันนี้ กรุณาเข้าสู่ระบบเพื่อค้นหาเพิ่มเติม",
					"action":  "login_required",
				},
			})
		},
	})
}

// MediaLimit limits media search for both guests and authenticated users
func (m *RateLimitMiddleware) MediaLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max: m.authConfig.MediaLimit,
		Expiration: 24 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			userID := getUserID(c)
			if userID != "" {
				return "auth:media:" + userID
			}
			return "guest:media:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			userID := getUserID(c)
			if userID != "" {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"success": false,
					"error": fiber.Map{
						"code":    "RATE_LIMITED",
						"message": "คุณค้นหารูปภาพ/วิดีโอเกินจำนวนที่กำหนดสำหรับวันนี้ (100 ครั้ง/วัน)",
						"action":  "wait",
					},
				})
			}
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "RATE_LIMITED",
					"message": "คุณค้นหารูปภาพ/วิดีโอเกินจำนวนที่กำหนดสำหรับวันนี้ กรุณาเข้าสู่ระบบเพื่อค้นหาเพิ่มเติม",
					"action":  "login_required",
				},
			})
		},
	})
}

// Legacy support - GuestRateLimitMiddleware for backward compatibility
type GuestRateLimitMiddleware struct {
	config GuestRateLimitConfig
}

func NewGuestRateLimitMiddleware(cfg GuestRateLimitConfig) *GuestRateLimitMiddleware {
	return &GuestRateLimitMiddleware{config: cfg}
}

// GuestSearchLimit limits search for non-logged users (legacy)
func (m *GuestRateLimitMiddleware) GuestSearchLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        m.config.SearchLimit,
		Expiration: 24 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			if c.Locals("user") != nil {
				return "user:" + c.IP()
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
		Next: func(c *fiber.Ctx) bool {
			return c.Locals("user") != nil
		},
	})
}

// GuestPlacesLimit limits places search for guests (legacy)
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

// GuestMediaLimit limits media search for guests (legacy)
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
