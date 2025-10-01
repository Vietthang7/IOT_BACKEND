package middleware

import (
	"backend/internal/auth"
	"backend/internal/consts"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Lấy token từ header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"code":    fiber.StatusUnauthorized,
				"error":   consts.GetFailed,
				"message": "Token không được cung cấp",
			})
		}
		// Check Bearer prefix
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"code":    fiber.StatusUnauthorized,
				"error":   consts.GetFailed,
				"message": "Format token không đúng",
			})
		}
		tokenString := tokenParts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"code":    fiber.StatusUnauthorized,
				"error":   consts.GetFailed,
				"message": "Token không hợp lệ hoặc đã hết hạn",
			})
		}
		c.Locals("username", claims.Username)
		return c.Next()
	}
}
