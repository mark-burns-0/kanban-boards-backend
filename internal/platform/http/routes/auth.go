package routes

import "github.com/gofiber/fiber/v2"

func AuthRoutes(router fiber.Router) fiber.Router {
	auth := router.Group("/auth")

	auth.Post("/login")    // вход в систему
	auth.Post("/register") // регистрация пользователя

	auth.Get("/current") // текущий пользователь
	auth.Put("/")        // обновление текущего пользователя

	return auth
}
