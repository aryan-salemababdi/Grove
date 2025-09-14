package app

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
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
	Middlewares() []string
}

type App struct {
	container *dig.Container
	http      *fiber.App
	modules   []string
	globalMws []string
}

func New() *App {
	return &App{
		container: dig.New(),
		http:      fiber.New(),
		modules:   []string{},
		globalMws: []string{},
	}
}

func (a *App) Container() *dig.Container { return a.container }
func (a *App) HTTP() *fiber.App          { return a.http }

func (a *App) UseGlobalMiddleware(names ...string) {
	for _, name := range names {
		if mw, ok := GetMiddleware(name); ok {
			a.http.Use(mw)
			a.globalMws = append(a.globalMws, name)
		} else {
			log.Println("⚠️ Global middleware not found:", name)
		}
	}
}

func (a *App) RegisterModule(name string, m Module) error {
	a.modules = append(a.modules, name)

	for _, mw := range m.Middlewares() {
		if fn, ok := GetMiddleware(mw); ok {
			a.http.Use(fn)
		} else {
			log.Println("⚠️ Module middleware not found:", mw)
		}
	}

	return m.Register(a.container, a.http)
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

func (a *App) Start(addr string) error {
	unregistered, err := a.checkUnregisteredModules("./")
	if err != nil {
		return fmt.Errorf("error scanning modules: %w", err)
	}
	if len(unregistered) > 0 {
		return errors.New("❌ Unregistered modules detected: " + strings.Join(unregistered, ", "))
	}

	fmt.Println("✅ Registered modules:", a.modules)
	fmt.Println("🌱 Velora starting on", addr)
	return a.http.Listen(addr)
}

func (a *App) WatchModules(rootDir string, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			unregistered, err := a.checkUnregisteredModules(rootDir)
			if err != nil {
				log.Println("Error scanning modules:", err)
				continue
			}
			if len(unregistered) > 0 {
				log.Println("❌ Unregistered modules detected:", strings.Join(unregistered, ", "))
			}
		}
	}()
}
