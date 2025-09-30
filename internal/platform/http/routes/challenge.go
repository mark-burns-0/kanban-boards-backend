package routes

import "github.com/gofiber/fiber/v2"

func ChallengeRoutes(router fiber.Router) fiber.Router {
	challenges := router.Group("/challenges")

	challenges.Get("/")               // получить список челленджей
	challenges.Get("/:uuid")          // получить челлендж по UUID
	challenges.Get("/:uuid/progress") // получить прогресс по челленджу у участников

	challenges.Post("/")        // создать новый челлендж
	challenges.Put("/:uuid")    // обновить челендж
	challenges.Delete("/:uuid") // удалить челендж

	return challenges
}
