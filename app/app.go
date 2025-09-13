package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"
)

type RouteRegisterer interface {
	RegisterRoutes(app *fiber.App)
}

type Module interface {
	Register(container *dig.Container, app *fiber.App) error
}

type App struct {
	container *dig.Container
	http      *fiber.App
}

func New() *App {
	return &App{
		container: dig.New(),
		http:      fiber.New(),
	}
}

func (a *App) Container() *dig.Container { return a.container }
func (a *App) HTTP() *fiber.App          { return a.http }

func (a *App) Provide(constructor interface{}) error {
	return a.container.Provide(constructor)
}

func (a *App) Invoke(fn interface{}) error {
	return a.container.Invoke(fn)
}

func (a *App) RegisterModule(m Module) error {
	return m.Register(a.container, a.http)
}

func (a *App) RegisterController(ctor interface{}) error {
	if err := a.container.Provide(ctor); err != nil {
		return err
	}
	return a.container.Invoke(func(r RouteRegisterer) {
		r.RegisterRoutes(a.http)
	})
}

func (a *App) Start(addr string) error {
	fmt.Println("Grove starting on", addr)
	return a.http.Listen(addr)
}
