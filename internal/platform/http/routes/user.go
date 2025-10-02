package routes

import (
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router fiber.Router, h AuthHandler) fiber.Router {
	users := router.Group("/users")

	users.Get("/current", h.Current)
	users.Put("/:id", h.Update)

	return users
}
