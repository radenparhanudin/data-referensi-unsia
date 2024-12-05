package routes

import "github.com/gofiber/fiber/v2"

func UserRoute(app fiber.Router) {
	user := app.Group("/users")
	user.Post("/login", func(c *fiber.Ctx) error {
		return c.SendString("Login")
	})
	user.Post("/logout", func(c *fiber.Ctx) error {
		return c.SendString("Logout")
	})
}
