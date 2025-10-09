package http

import (
	"backend/internal/platform/http/middleware"
	"backend/internal/platform/http/routes"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(r fiber.Router, handlers Handlers) {
	r.Use(middleware.Recover)
	r.Use(middleware.LogRequest)
	api := r.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/health", healthCheck)

	routes.AuthRoutes(v1, handlers.AuthHandler)
	routes.UserRoutes(v1, handlers.UserHandler)
	routes.BoardRoutes(v1, handlers.BoardHandler)
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(Response{
		Message:    "Ok",
		StatusCode: http.StatusOK,
	})
}
