package http

import (
	"github.com/gofiber/fiber/v2"
)

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Challenge Tracker",
	})
	return app
}
