package main

import (
	"fiber/config"
	"fiber/pkg/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.ConnectDB()
	app := fiber.New()
	routes.Route(app)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello,world")
	})

	app.Listen(":3000")
}
