package routes

import (
	"github.com/gofiber/fiber/v2"
)

func BiodataRoute(app fiber.Router) {
	biodata := app.Group("/biodatas")
	biodata.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "List biodata",
		})
	})
}
