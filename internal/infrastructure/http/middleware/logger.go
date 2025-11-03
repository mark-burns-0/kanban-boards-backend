package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LogRequest(c *fiber.Ctx) error {
	start := time.Now()

	c.Next()

	slog.Info("Request completed",
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
		slog.String("ip", c.IP()),
		slog.String("user-agent", string(c.Request().Header.UserAgent())),
		slog.Int("status-code", c.Response().StatusCode()),
		slog.Duration("duration", time.Since(start)),
	)

	return nil
}
