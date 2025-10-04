package validation

import (
	"backend/internal/platform/lang"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Validator struct {
	validator    *validator.Validate
	langRegistry *lang.Registry
}

func New() *Validator {
	return &Validator{
		validator:    validator.New(),
		langRegistry: lang.NewRegistry(),
	}
}

func (v *Validator) ValidateStruct(c *fiber.Ctx, structPtr any) error {
	err := v.validator.Struct(structPtr)
	if err != nil {
		validationMessages, err := v.langRegistry.Validate(c.Context(), err)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusUnprocessableEntity).
			JSON(fiber.Map{"error": validationMessages})
	}
	return nil
}
