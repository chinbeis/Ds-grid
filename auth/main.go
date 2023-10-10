package main

import (
	routes "auth/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	routes.R_login(app.Group("/login"))
	routes.R_users(app.Group("/auth/user/3"))

	app.Listen(":8080")
}
