package routes

import (
	"fiber/app/controllers"
	"fiber/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	route := app.Group("/api/v1/user")
	route.Post("/", controllers.SignupUser)
	route.Post("/login", controllers.LoginUser)
	route.Get("/", middleware.AuthMiddleware(), controllers.GetUserDetail)
	route.Get("/logout", middleware.AuthMiddleware(), controllers.Logout)
	route.Get("/refresh", controllers.Refresh)
}
