package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func LogRequest(c *fiber.Ctx) error {
	slog.Info("Request received",
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
		slog.String("ip", c.IP()),
		slog.String("user-agent", string(c.Request().Header.UserAgent())),
	)
	return c.Next()
}
