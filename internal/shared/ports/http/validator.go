package http

import "github.com/gofiber/fiber/v2"

type Validator interface {
	ValidateStruct(c *fiber.Ctx, structPtr interface{}) (map[string]string, int, error)
}
