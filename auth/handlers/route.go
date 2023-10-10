package handlers

import (
	"auth/service"

	"github.com/gofiber/fiber/v2"
)

func R_login(router fiber.Router) {
	router.Get("/", func(context *fiber.Ctx) error {
		return context.SendString("respond with a resource")
	})
}

func R_users(router fiber.Router) {
	router.Get("/", func(context *fiber.Ctx) error {
		return service.DsGrid(context)
	})
}
