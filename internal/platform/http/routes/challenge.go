package routes

import (
	"backend/internal/platform/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type ChallengeHandler interface {
	GetList(*fiber.Ctx) error
	GetByUUID(*fiber.Ctx) error
	GetProgress(*fiber.Ctx) error
	Store(*fiber.Ctx) error
	Update(*fiber.Ctx) error
	Delete(*fiber.Ctx) error
}

func ChallengeRoutes(router fiber.Router, handler ChallengeHandler) fiber.Router {
	challenges := router.Group("/challenges").Use(middleware.AuthRequired)

	challenges.Get("/", handler.GetList)                   // получить список челленджей
	challenges.Get("/:uuid", handler.GetByUUID)            // получить челлендж по UUID
	challenges.Get("/:uuid/progress", handler.GetProgress) // получить прогресс по челленджу у участников

	challenges.Post("/", handler.Store)         // создать новый челлендж
	challenges.Put("/:uuid", handler.Update)    // обновить челендж
	challenges.Delete("/:uuid", handler.Delete) // удалить челендж

	return challenges
}
