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
}

func BoardRoutes(router fiber.Router, handler BoardHandler) fiber.Router {
	boards := router.Group("/boards").Use(middleware.AuthRequired)

	boards.Get("/", handler.GetList)      // получить список челленджей
	boards.Get("/:id", handler.GetByUUID) // получить челлендж по UUID

	boards.Post("/", handler.Store)       // создать новый челлендж
	boards.Put("/:id", handler.Update)    // обновить челендж
	boards.Delete("/:id", handler.Delete) // удалить челендж

	return boards
}
