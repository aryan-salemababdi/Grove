package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "g",
	Short: "generate resources (e.g. modules)",
}

var genModuleCmd = &cobra.Command{
	Use:   "module [name]",
	Short: "generate a new module scaffold",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return generateModule(name)
	},
}

func init() {
	genCmd.AddCommand(genModuleCmd)
}

// templates
var serviceTmpl = `package {{.Pkg}}

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) FindAll() []string { return []string{"example"} }
`

var controllerTmpl = `package {{.Pkg}}

import "github.com/gofiber/fiber/v2"

type Controller struct {
	Service *Service
}

func NewController(s *Service) *Controller { return &Controller{Service: s} }

func (c *Controller) RegisterRoutes(app *fiber.App) {
	app.Get("/{{.Route}}", func(ctx *fiber.Ctx) error { return ctx.JSON(c.Service.FindAll()) })
}
`

var moduleTmpl = `package {{.Pkg}}

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"
)

type {{.UpName}}Module struct{}

func New() *{{.UpName}}Module { return &{{.UpName}}Module{} }

func (m *{{.UpName}}Module) Register(container *dig.Container, app *fiber.App) error {
	if err := container.Provide(NewService); err != nil { return err }
	if err := container.Provide(NewController); err != nil { return err }
	return container.Invoke(func(ctrl *Controller) { ctrl.RegisterRoutes(app) })
}
`

func generateModule(name string) error {
	pkg := name
	up := string(bytes.ToUpper([]byte(name[:1]))) + name[1:]
	root := filepath.Join(".", name)
	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}
	files := map[string]string{
		"service.go":    serviceTmpl,
		"controller.go": controllerTmpl,
		"module.go":     moduleTmpl,
	}
	for fname, tmpl := range files {
		path := filepath.Join(root, fname)
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		t := template.Must(template.New(fname).Parse(tmpl))
		data := map[string]string{"Pkg": pkg, "UpName": up, "Route": pkg}
		if err := t.Execute(f, data); err != nil {
			return err
		}
		fmt.Println("created", path)
	}
	return nil
}
