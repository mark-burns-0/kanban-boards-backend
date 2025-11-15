package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func CORS(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
	c.Set("Access-Control-Allow-Headers", "Content-Type,Authorization,Refresh-Token,Accept,Origin,X-Requested-With")
	c.Set("Access-Control-Allow-Credentials", "true")
	c.Set("Access-Control-Max-Age", "3600")

	if c.Method() == "OPTIONS" {
		return c.SendStatus(fiber.StatusOK)
	}

	return c.Next()
}
