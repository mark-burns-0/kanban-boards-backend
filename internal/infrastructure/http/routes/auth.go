package routes

import (
	"backend/internal/infrastructure/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(*fiber.Ctx) error
	Register(*fiber.Ctx) error
	RefreshToken(*fiber.Ctx) error
}

func AuthRoutes(router fiber.Router, h AuthHandler) fiber.Router {
	auth := router.Group("/auth")

	auth.Post("/login", h.Login)
	auth.Post("/register", h.Register)
	auth.Post("/refresh-token", middleware.RefreshToken, h.RefreshToken)

	return auth
}
