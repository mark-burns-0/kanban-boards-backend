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

	auth.Post("/login", h.Login)       // вход в систему
	auth.Post("/register", h.Register) // регистрация пользователя
	auth.Post("/refresh-token", middleware.AuthRequired, h.RefreshToken)

	return auth
}
