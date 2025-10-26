package utils

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

const (
	UserIDKey   = "userID"
	FieldUserID = "UserID"
)

func ParseBody[T any](c *fiber.Ctx) (*T, error) {
	var body T
	if err := c.BodyParser(&body); err != nil {
		return nil, err
	}

	setUserID(&body, c)

	return &body, nil
}

func setUserID(body any, c *fiber.Ctx) {
	val := reflect.ValueOf(body).Elem()
	field := val.FieldByName(FieldUserID)
	if !field.IsValid() || !field.CanSet() {
		return
	}

	if userID, ok := c.Locals(UserIDKey).(uint64); ok && field.Kind() == reflect.Uint64 {
		field.SetUint(userID)
	}
}
