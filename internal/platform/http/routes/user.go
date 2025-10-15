package routes

import (
	"backend/internal/platform/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	Current(*fiber.Ctx) error
	Update(*fiber.Ctx) error
}

func UserRoutes(router fiber.Router, h UserHandler) fiber.Router {
	users := router.Group("/users").Use(middleware.AuthRequired)

	users.Get("/current", h.Current)
	users.Put("/current", h.Update)

	return users
}
