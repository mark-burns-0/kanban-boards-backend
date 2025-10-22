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
	r.Use(middleware.MethodWhiteList)
	api := r.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/health", healthCheck)

	routes.AuthRoutes(v1, handlers.AuthHandler)
	routes.UserRoutes(v1, handlers.UserHandler)
	routes.BoardRoutes(v1, handlers.BoardHandler)
	routes.CardRoutes(v1, handlers.CardHandler)
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(Response{
		Message:    "OK",
		StatusCode: http.StatusOK,
	})
}
