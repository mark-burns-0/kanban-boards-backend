package routes

import (
	"backend/internal/infrastructure/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type CardHandler interface {
	Create(*fiber.Ctx) error
	Delete(*fiber.Ctx) error
	Update(*fiber.Ctx) error
	MoveToNewPosition(*fiber.Ctx) error
}

func CardRoutes(router fiber.Router, h CardHandler) fiber.Router {
	cards := router.Group("/boards/:id/cards").
		Use(middleware.AuthRequired)

	cards.Post("/", h.Create)

	cardIDGroup := cards.Group("/:card_id")
	cardIDGroup.Delete("/", h.Delete)
	cardIDGroup.Put("/", h.Update)
	cardIDGroup.Put("/move", h.MoveToNewPosition)

	return cards
}
