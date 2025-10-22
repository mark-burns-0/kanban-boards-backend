package routes

import (
	"backend/internal/platform/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type CommentHandler interface {
	Create(*fiber.Ctx) error
	Update(*fiber.Ctx) error
	Delete(*fiber.Ctx) error
}

func CommentRoutes(router fiber.Router, h CommentHandler) fiber.Router {
	comments := router.Group("boards/:id/cards/:card_id/comments").
		Use(middleware.AuthRequired)

	comments.Post("/", h.Create)
	comments.Put("/:comment_id", h.Update)
	comments.Delete("/:comment_id", h.Delete)

	return comments
}
