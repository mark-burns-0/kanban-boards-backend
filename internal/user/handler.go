package user

import "github.com/gofiber/fiber/v2"

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error    { return nil }
func (h *AuthHandler) Register(c *fiber.Ctx) error { return nil }
func (h *AuthHandler) Current(c *fiber.Ctx) error  { return nil }
func (h *AuthHandler) Update(c *fiber.Ctx) error   { return nil }
