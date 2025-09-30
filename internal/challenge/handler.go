package challenge

import "github.com/gofiber/fiber/v2"

type ChallengeHandler struct{}

func NewChallengeHandler() *ChallengeHandler {
	return &ChallengeHandler{}
}

func (h *ChallengeHandler) GetByUUID(c *fiber.Ctx) error   { return nil }
func (h *ChallengeHandler) GetList(c *fiber.Ctx) error     { return nil }
func (h *ChallengeHandler) GetProgress(c *fiber.Ctx) error { return nil }
func (h *ChallengeHandler) Store(c *fiber.Ctx) error       { return nil }
func (h *ChallengeHandler) Update(c *fiber.Ctx) error      { return nil }
func (h *ChallengeHandler) Delete(c *fiber.Ctx) error      { return nil }
