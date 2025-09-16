package app

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HttpException struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	ErrorText  string      `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func (e *HttpException) Error() string {
	return e.Message
}

func NewHttpException(statusCode int, message string, data interface{}) *HttpException {
	return &HttpException{
		StatusCode: statusCode,
		Message:    message,
		ErrorText:  http.StatusText(statusCode),
		Data:       data,
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	if e, ok := err.(*HttpException); ok {
		return c.Status(e.StatusCode).JSON(e)
	}

	if e, ok := err.(*fiber.Error); ok {
		return c.Status(e.Code).JSON(HttpException{
			StatusCode: e.Code,
			Message:    e.Message,
			ErrorText:  http.StatusText(e.Code),
		})
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(HttpException{
			StatusCode: fiber.StatusNotFound,
			Message:    "Record not found",
			ErrorText:  "Not Found",
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(HttpException{
		StatusCode: fiber.StatusInternalServerError,
		Message:    err.Error(),
		ErrorText:  "Internal Server Error",
	})
}
