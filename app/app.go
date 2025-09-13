package app

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

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
	modules   []string
}

func New() *App {
	return &App{
		container: dig.New(),
		http:      fiber.New(),
		modules:   []string{},
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

func (a *App) RegisterModule(name string, m Module) error {
	a.modules = append(a.modules, name)
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

func (a *App) checkUnregisteredModules(rootDir string) ([]string, error) {
	unregistered := []string{}
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".module.go") {
			modName := strings.TrimSuffix(d.Name(), ".module.go")
			registered := false
			for _, r := range a.modules {
				if r == modName {
					registered = true
					break
				}
			}
			if !registered {
				unregistered = append(unregistered, modName)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return unregistered, nil
}

func (a *App) WatchModules(rootDir string, interval time.Duration) {
	go func() {
		for {
			unregistered, err := a.checkUnregisteredModules(rootDir)
			if err != nil {
				fmt.Println("Error scanning modules:", err)
			} else if len(unregistered) > 0 {
				fmt.Println("❌ Unregistered modules detected:", strings.Join(unregistered, ", "))
			}
			time.Sleep(interval)
		}
	}()
}

func (a *App) Start(addr string) error {

	unregistered, err := a.checkUnregisteredModules("./")
	if err != nil {
		return fmt.Errorf("error scanning modules: %w", err)
	}
	if len(unregistered) > 0 {
		return errors.New("❌ The following modules are not registered in AppModule: " + strings.Join(unregistered, ", "))
	}

	fmt.Println("✅ Registered modules:", a.modules)
	fmt.Println("🌱 Grove starting on", addr)
	return a.http.Listen(addr)
}
