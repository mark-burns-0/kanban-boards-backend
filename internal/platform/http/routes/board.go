package routes

import (
	"backend/internal/platform/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type BoardHandler interface {
	GetList(*fiber.Ctx) error
	GetByUUID(*fiber.Ctx) error
	Store(*fiber.Ctx) error
	Update(*fiber.Ctx) error
	Delete(*fiber.Ctx) error
	MoveToColumn(*fiber.Ctx) error
}

func BoardRoutes(router fiber.Router, handler BoardHandler) fiber.Router {
	boards := router.Group("/boards").Use(middleware.AuthRequired)

	boards.Get("/", handler.GetList)      // получить список досок
	boards.Get("/:id", handler.GetByUUID) // получить доску по UUID

	boards.Post("/", handler.Store)       // создать новую доску
	boards.Put("/:id", handler.Update)    // обновить доску
	boards.Delete("/:id", handler.Delete) // удалить доску

	boards.Put("/:id/columns/:column_id/move", handler.MoveToColumn) // переместить в колонку

	return boards
}
