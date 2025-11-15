package middleware

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func Recover(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered",
				slog.String("errors", fmt.Sprintf("%v", r)),
				slog.String("path", c.Path()),
				slog.String("method", c.Method()),
			)
			// Отправляем ошибку 500 клиенту
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"errors": "Internal Server Error",
			})
		}
	}()

	return c.Next()
}
