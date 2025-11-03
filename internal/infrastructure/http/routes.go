package http

import (
	"backend/internal/infrastructure/http/middleware"
	"backend/internal/infrastructure/http/routes"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(r fiber.Router, handlers Handlers) {
	r.Use(middleware.LogRequest)
	r.Use(middleware.Recover)
	r.Use(middleware.NotFound)
	r.Use(middleware.MethodWhiteList)
	api := r.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/health", healthCheck)

	routes.AuthRoutes(v1, handlers.AuthHandler)
	routes.UserRoutes(v1, handlers.UserHandler)
	routes.BoardRoutes(v1, handlers.BoardHandler)
	routes.CardRoutes(v1, handlers.CardHandler)
	routes.CommentRoutes(v1, handlers.CommentHandler)
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(struct {
		Message    string `json:"message"`
		StatusCode int    `json:"status_code"`
	}{
		StatusCode: 200,
		Message:    "Ok",
	},
	)
}
