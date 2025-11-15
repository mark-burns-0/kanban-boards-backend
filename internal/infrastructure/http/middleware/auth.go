package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	UserIDKey    = "userID"
	BearerPrefix = "Bearer "
)

type sub struct {
	UserID *uint64 `json:"id"`
}

type Claims struct {
	Sub sub `json:"sub"`
	jwt.RegisteredClaims
}

func AuthRequired(c *fiber.Ctx) error {
	token := getToken(c)

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(
				fiber.Map{"errors": "missing_token"},
			)
	}

	claims := &Claims{}
	if parsedToken, err := parseToken(token, claims); err != nil || !parsedToken.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{"errors": "token_expired_or_invalid"},
		)
	}

	if claims.Sub.UserID == nil || *claims.Sub.UserID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{"errors": "invalid_user_identifier"},
		)
	}

	c.Locals(UserIDKey, *claims.Sub.UserID)

	return c.Next()
}

func getToken(c *fiber.Ctx) string {
	token := c.Get("Authorization")
	token = strings.TrimPrefix(token, BearerPrefix)
	return token
}

func parseToken(token string, claims *Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_ACCESS_TOKEN_SECRET")), nil
	})
}
