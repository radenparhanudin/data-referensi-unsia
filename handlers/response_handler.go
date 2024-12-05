package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func SendSuccess(c *fiber.Ctx, statusCode int, data interface{}, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"error":   false,
		"data":    data,
		"message": message,
	})
}

func SendFailed(c *fiber.Ctx, statusCode int, data interface{}, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"error":   true,
		"data":    nil,
		"message": message,
	})
}
