package handler

import "github.com/gofiber/fiber/v2"

func ResponseSuccess(c *fiber.Ctx, code int, message string, data interface{}) error {
	response := fiber.Map{
		"status":  true,
		"code":    code,
		"data":    data,
		"message": message,
	}
	return c.Status(code).JSON(response)
}

func ResponseError(c *fiber.Ctx, code int, message interface{}, err interface{}) error {
	return c.Status(code).JSON(fiber.Map{
		"status":  false,
		"code":    code,
		"error":   err,
		"message": message,
	})
}
