package middleware

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func Recover(c *fiber.Ctx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// Логируем ошибку. Здесь можно привести r к ошибке или строке
			slog.Error("Panic recovered", slog.String("error", fmt.Sprintf("%v", r)))
			// Возвращаем ошибку, чтобы Fiber обработал её
			err = fiber.ErrInternalServerError
		}
	}()
	return c.Next()
}
