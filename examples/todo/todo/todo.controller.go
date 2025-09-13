package todo

import "github.com/gofiber/fiber/v2"

type Controller struct {
	Service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{Service: s}
}

func (c *Controller) RegisterRoutes(app *fiber.App) {
	app.Get("/todos", func(ctx *fiber.Ctx) error {
		return ctx.JSON(c.Service.FindAll())
	})
}
