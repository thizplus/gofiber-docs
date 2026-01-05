package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// RequireLogin middleware - blocks non-logged users with login prompt response
func RequireLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("user") == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "LOGIN_REQUIRED",
					"message": "กรุณาเข้าสู่ระบบเพื่อใช้งานฟีเจอร์นี้",
					"action":  "login_required",
					"feature": c.Path(),
				},
			})
		}
		return c.Next()
	}
}
