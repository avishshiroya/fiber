package routes

import (
	"fiber/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	route := app.Group("/api/v1/user")
	route.Post("/",controllers.SignupUser)
	route.Post("/login",controllers.LoginUser)
}
