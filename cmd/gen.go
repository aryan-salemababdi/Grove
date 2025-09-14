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
		fmt.Sprintf("%s.service.go", name):    serviceTmpl,
		fmt.Sprintf("%s.controller.go", name): controllerTmpl,
		fmt.Sprintf("%s.module.go", name):     moduleTmpl,
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

	testDir := filepath.Join(root, "__tests__")
	if err := os.Mkdir(testDir, 0755); err != nil {
		return err
	}

	serviceTestTmpl := `package {{.Pkg}}

import "testing"

func TestServiceFindAll(t *testing.T) {
	s := NewService()
	got := s.FindAll()
	want := []string{"example"}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
`

	controllerTestTmpl := `package {{.Pkg}}

import (
	"testing"
	"github.com/gofiber/fiber/v2"
	"net/http/httptest"
)

func TestControllerRoute(t *testing.T) {
	s := NewService()
	ctrl := NewController(s)
	app := fiber.New()
	ctrl.RegisterRoutes(app)
	req := httptest.NewRequest("GET", "/{{.Route}}", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200; got %d", resp.StatusCode)
	}
}
`

	testFiles := map[string]string{
		fmt.Sprintf("%s.serviceـtest.go", name):    serviceTestTmpl,
		fmt.Sprintf("%s.controllerـtest.go", name): controllerTestTmpl,
	}

	for fname, tmpl := range testFiles {
		path := filepath.Join(testDir, fname)
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		t := template.Must(template.New(fname).Parse(tmpl))
		data := map[string]string{"Pkg": pkg, "Route": pkg}
		if err := t.Execute(f, data); err != nil {
			return err
		}
		fmt.Println("created", path)
	}

	return nil
}
