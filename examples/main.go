package main

import (
	"github.com/aryan-salemababdi/Grove/examples/user"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	userService := &user.UserService{}
	userController := &user.Controller{Service: userService}
	userController.RegisterRoutes(app)

	app.Listen(":4000")
}
