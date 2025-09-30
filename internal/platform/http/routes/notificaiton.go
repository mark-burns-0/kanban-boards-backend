package routes

import (
	"backend/internal/platform/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler interface {
	GetUnread(*fiber.Ctx) error
	MarkAsRead(*fiber.Ctx) error
}

func NotificationRoutes(router fiber.Router, handler NotificationHandler) fiber.Router {
	notifications := router.Group("/notifications").Use(middleware.AuthRequired)

	notifications.Get("/", handler.GetUnread)
	notifications.Post("/:uuid/read", handler.MarkAsRead)

	return notifications
}
