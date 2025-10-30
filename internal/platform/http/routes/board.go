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
	CreateColumn(*fiber.Ctx) error
	UpdateColumn(*fiber.Ctx) error
	DeleteColumn(*fiber.Ctx) error
	MoveColumn(*fiber.Ctx) error
}

func BoardRoutes(router fiber.Router, handler BoardHandler) fiber.Router {
	boards := router.Group("/boards").Use(middleware.AuthRequired)

	boards.Get("/:id", handler.GetByUUID) // получить доску по UUID

	boards.Post("/list", handler.GetList) // получить список досок
	boards.Post("/", handler.Store)       // создать новую доску
	boards.Put("/:id", handler.Update)    // обновить доску
	boards.Delete("/:id", handler.Delete) // удалить доску

	// Работа с колонками
	columns := boards.Group("/:id/columns")

	columns.Post("/", handler.CreateColumn)
	columns.Put("/:column_id", handler.UpdateColumn)
	columns.Put("/:column_id/move", handler.MoveColumn)
	columns.Delete("/:column_id", handler.DeleteColumn)

	return boards
}
