package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func MethodWhiteList(c *fiber.Ctx) error {
	allowedMethods := map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"OPTIONS": true,
	}

	if !allowedMethods[c.Method()] {
		return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
			"errors": fmt.Sprintf("Method %s not allowed", c.Method()),
		})
	}

	return c.Next()
}
