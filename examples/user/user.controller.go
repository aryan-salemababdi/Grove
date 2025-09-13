package user

import "github.com/gofiber/fiber/v2"

type Controller struct {
	Service *UserService
}

func NewController(s *UserService) *Controller {
	return &Controller{Service: s}
}

func (c *Controller) RegisterRoutes(app *fiber.App) {
	app.Get("/users", func(ctx *fiber.Ctx) error {
		return ctx.JSON(c.Service.FindAll())
	})
}
