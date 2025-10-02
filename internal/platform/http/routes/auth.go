package routes

import (
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(*fiber.Ctx) error
	Register(*fiber.Ctx) error
	Current(*fiber.Ctx) error
	Update(*fiber.Ctx) error
	Refresh(*fiber.Ctx) error
}

func AuthRoutes(router fiber.Router, h AuthHandler) fiber.Router {
	auth := router.Group("/auth")

	auth.Post("/login", h.Login)       // вход в систему
	auth.Post("/register", h.Register) // регистрация пользователя

	return auth
}
