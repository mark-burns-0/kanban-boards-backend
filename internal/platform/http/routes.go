package http

import (
	"backend/internal/platform/http/routes"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(r fiber.Router) {
	api := r.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/health", healthCheck)

	routes.AuthRoutes(v1)
	routes.ChallengeRoutes(v1)
	routes.NotificationRoutes(v1)
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(Response{
		Message:    "Ok",
		StatusCode: http.StatusOK,
	})
}
