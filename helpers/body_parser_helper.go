package helpers

import "github.com/gofiber/fiber/v2"

func BodyParser(c *fiber.Ctx, body interface{}) error {
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid data format",
		})
	}
	return nil
}
