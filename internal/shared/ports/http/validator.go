package http

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type Validator interface {
	ValidateStruct(c *fiber.Ctx, structPtr any) (map[string]string, int, error)
}

type LangMessage interface {
	GetResponseMessage(ctx context.Context, key string) string
}
