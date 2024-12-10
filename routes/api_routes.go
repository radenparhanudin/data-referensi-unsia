package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
	api := app.Group("/api")
	RegionRoute(api)
	BiodataRoute(api)
	EducationRoute(api)
}
