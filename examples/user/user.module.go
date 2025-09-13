package user

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"
)

type UserModule struct{}

func New() *UserModule { return &UserModule{} }

func (m *UserModule) Register(container *dig.Container, app *fiber.App) error {

	if err := container.Provide(NewService); err != nil {
		return err
	}
	if err := container.Provide(NewController); err != nil {
		return err
	}

	return container.Invoke(func(ctrl *Controller) {
		ctrl.RegisterRoutes(app)
	})
}
