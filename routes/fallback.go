package routes

import "github.com/gofiber/fiber/v2"

func FallbackRoute(app *fiber.App) {
	app.All("*", func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"message": "The requested route does not exist.",
		})
	})
}
