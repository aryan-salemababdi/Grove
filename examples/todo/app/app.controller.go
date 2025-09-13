package app

import "github.com/gofiber/fiber/v2"

type Controller struct {
	Service *Service
}

func NewController(s *Service) *Controller { return &Controller{Service: s} }

func (c *Controller) RegisterRoutes(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString(c.Service.Greet())
	})
}
