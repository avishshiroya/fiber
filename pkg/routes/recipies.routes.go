package routes

import (
	"fiber/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func RecipieRoute(app *fiber.App) {
	route := app.Group("/api/v1/recipies")
	route.Post("/", controllers.CreateRecipies)
	route.Post("/notification", controllers.CreateNotification)
}
