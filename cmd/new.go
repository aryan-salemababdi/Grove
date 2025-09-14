package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new Grove application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return createNewApp(name)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func createNewApp(name string) error {
	if err := os.Mkdir(name, 0755); err != nil {
		return err
	}

	modInit := exec.Command("go", "mod", "init", name)
	modInit.Dir = name
	if out, err := modInit.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod init failed: %s", string(out))
	}

	mainContent := `package main

import (
	"log"
	"time"
	"fmt"

	grove "github.com/aryan-salemababdi/Grove/app"
	"{{.Name}}/app"
)

func main() {
	a := grove.New()

    grove.RegisterMiddleware("auth", func(c *fiber.Ctx) error {
        return c.Next()
    })

    grove.RegisterMiddleware("logging", func(c *fiber.Ctx) error {
        fmt.Println("[LOG]", c.Path())
        return c.Next()
    })

    grove.RegisterMiddleware("cors", func(c *fiber.Ctx) error {
        c.Set("Access-Control-Allow-Origin", "*")
        return c.Next()
    })

    a.UseGlobalMiddleware("logging", "cors", "auth")

	if err := a.RegisterModule("app", app.New()); err != nil {
		log.Fatal(err)
	}

	a.WatchModules(".", 2*time.Second)

	if err := a.Start(":3000"); err != nil {
		log.Fatal(err)
	}
}
`
	if err := writeTemplate(filepath.Join(name, "main.go"), mainContent, map[string]string{"Name": name}); err != nil {
		return err
	}

	moduleDir := filepath.Join(name, "app")
	if err := os.Mkdir(moduleDir, 0755); err != nil {
		return err
	}

	// Service
	service := `package app

type Service struct{}

func NewService() *Service  { return &Service{} }

func (s *Service) Greet() string { return  "hello from Grove!" }
`
	if err := writeFile(filepath.Join(moduleDir, "app.service.go"), service); err != nil {
		return err
	}

	// Controller
	controller := `package app

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
`
	if err := writeFile(filepath.Join(moduleDir, "app.controller.go"), controller); err != nil {
		return err
	}

	// Module
	module := `package app

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"
)

type AppModule struct{}

func New() *AppModule { return &AppModule{} }

func (m *AppModule) Register(container *dig.Container, app *fiber.App) error {
	if err := container.Provide(NewService); err != nil { return err }
	if err := container.Provide(NewController); err != nil { return err }
	return container.Invoke(func(ctrl *Controller) { ctrl.RegisterRoutes(app) })
}

func (m *AppModule) Middlewares() []string {
	return []string{}
}
`
	if err := writeFile(filepath.Join(moduleDir, "app.module.go"), module); err != nil {
		return err
	}

	modTidy := exec.Command("go", "mod", "tidy")
	modTidy.Dir = name
	if out, err := modTidy.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy failed: %s", string(out))
	}

	fmt.Println("✅ New Grove app created at", name)
	fmt.Println("👉 cd", name, "&& go run main.go")
	return nil
}

func writeFile(path, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

func writeTemplate(path, tmpl string, data map[string]string) error {
	t := template.Must(template.New("file").Parse(tmpl))
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, data)
}
