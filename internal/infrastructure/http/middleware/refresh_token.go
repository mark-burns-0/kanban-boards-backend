package middleware

import "github.com/gofiber/fiber/v2"

var (
	RefreshTokenKey    = "refreshToken"
	RefreshTokenHeader = "Refresh-Token"
)

func RefreshToken(c *fiber.Ctx) error {
	refreshToken := getRefreshToken(c)
	if refreshToken == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	setRefreshToken(c, refreshToken)
	return c.Next()
}

func getRefreshToken(c *fiber.Ctx) string {
	return c.Get(RefreshTokenHeader)
}

func setRefreshToken(c *fiber.Ctx, token string) {
	c.Locals(RefreshTokenKey, token)
}
