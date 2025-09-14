package testing

import (
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"

	"github.com/aryan-salemababdi/Velora/app"
)

func SetupTestApp(modules ...app.Module) (*fiber.App, *dig.Container) {
	a := app.New()
	for _, m := range modules {
		if err := a.RegisterModule("test_module", m); err != nil {
			panic(err)
		}
	}
	return a.HTTP(), a.Container()
}

func TestRequest(app *fiber.App, method, path string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	resp := httptest.NewRecorder()
	app.Test(req)
	return resp
}
