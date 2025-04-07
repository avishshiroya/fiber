package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	UserRoutes(app)
	RecipieRoute(app)
}
