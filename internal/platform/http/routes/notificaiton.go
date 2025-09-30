package routes

import "github.com/gofiber/fiber/v2"

func NotificationRoutes(router fiber.Router) fiber.Router {
	notifications := router.Group("/notifications")

	notifications.Get("/")
	notifications.Post("/:uuid/read")

	return notifications
}
