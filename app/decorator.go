package app

import "github.com/gofiber/fiber/v2"

type Decorator func(fiber.Handler) fiber.Handler

func Use(decorators ...Decorator) func(fiber.Handler) fiber.Handler {
	return func(h fiber.Handler) fiber.Handler {
		for i := len(decorators) - 1; i >= 0; i-- {
			h = decorators[i](h)
		}
		return h
	}
}
