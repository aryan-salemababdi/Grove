package app

import (
	"todo/todo"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"
)

type AppModule struct{}

func New() *AppModule { return &AppModule{} }

func (m *AppModule) Register(container *dig.Container, app *fiber.App) error {
	if err := container.Provide(NewService); err != nil {
		return err
	}
	if err := container.Provide(NewController); err != nil {
		return err
	}
	if err := container.Invoke(func(ctrl *Controller) { ctrl.RegisterRoutes(app) }); err != nil {
		return err
	}

	return todo.New().Register(container, app)
}
