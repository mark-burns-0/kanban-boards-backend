package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	timeout = 5 * time.Second
)

type Config interface {
	GetAppName() string
}

func NewApp(cfg Config) *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:      true,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		AppName:      cfg.GetAppName(),
	})
	return app
}
