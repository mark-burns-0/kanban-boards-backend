package http

import (
	"github.com/gofiber/fiber/v2"
)

type Config interface {
	GetAppName() string
}

func NewApp(cfg Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: cfg.GetAppName(),
	})
	return app
}
