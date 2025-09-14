package app

import "github.com/gofiber/fiber/v2"

type MiddlewareFunc = fiber.Handler

var middlewareRegistry = map[string]MiddlewareFunc{}

func RegisterMiddleware(name string, fn MiddlewareFunc) {
	middlewareRegistry[name] = fn
}

func GetMiddleware(name string) (MiddlewareFunc, bool) {
	fn, ok := middlewareRegistry[name]
	return fn, ok
}
