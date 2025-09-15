package app

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func AutoBind(dtoType interface{}) Decorator {
	return func(next fiber.Handler) fiber.Handler {
		return func(c *fiber.Ctx) error {
			dto := reflect.New(reflect.TypeOf(dtoType).Elem()).Interface()

			if err := c.BodyParser(dto); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "invalid request body",
				})
			}

			if err := validate.Struct(dto); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			c.Locals("dto", dto)

			return next(c)
		}
	}
}
