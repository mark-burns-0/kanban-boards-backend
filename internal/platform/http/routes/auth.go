package routes

import (
	"backend/internal/platform/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(*fiber.Ctx) error
	Register(*fiber.Ctx) error
	Current(*fiber.Ctx) error
	Update(*fiber.Ctx) error
}

func AuthRoutes(router fiber.Router, h AuthHandler) fiber.Router {
	auth := router.Group("/auth")

	auth.Post("/login", h.Login)       // вход в систему
	auth.Post("/register", h.Register) // регистрация пользователя

	auth.Get("/current", middleware.AuthRequired, h.Current) // текущий пользователь
	auth.Put("/", middleware.AuthRequired, h.Update)         // обновление текущего пользователя

	return auth
}
