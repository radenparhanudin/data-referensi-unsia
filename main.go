package main

import (
	"data-referensi/config"
	"data-referensi/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	config.ConnectDB()

	routes.SetupRouter(app)

	routes.FallbackRoute(app)

	app.Listen(":3000")
}