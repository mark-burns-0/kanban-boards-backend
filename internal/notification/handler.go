package notification

import "github.com/gofiber/fiber/v2"

type NotificationHandler struct{}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (h *NotificationHandler) GetUnread(c *fiber.Ctx) error  { return nil }
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error { return nil }
