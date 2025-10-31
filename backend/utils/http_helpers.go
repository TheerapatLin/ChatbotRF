package utils

import (
	"github.com/gofiber/fiber/v2"
)

// ErrorResponse returns a JSON error response with the specified status code
func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}

// BadRequest returns a 400 Bad Request error
func BadRequest(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message)
}

// NotFound returns a 404 Not Found error
func NotFound(c *fiber.Ctx, resource string, id interface{}) error {
	message := "Resource not found"
	if resource != "" && id != nil {
		message = resource + " not found"
	}
	return ErrorResponse(c, fiber.StatusNotFound, message)
}

// InternalError returns a 500 Internal Server Error
func InternalError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, message)
}

// ServiceUnavailable returns a 503 Service Unavailable error
func ServiceUnavailable(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusServiceUnavailable, message)
}

// UnprocessableEntity returns a 422 Unprocessable Entity error
func UnprocessableEntity(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnprocessableEntity, message)
}

// UnsupportedMediaType returns a 415 Unsupported Media Type error
func UnsupportedMediaType(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnsupportedMediaType, message)
}

// RequestEntityTooLarge returns a 413 Request Entity Too Large error
func RequestEntityTooLarge(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusRequestEntityTooLarge, message)
}

// SuccessJSON returns a 200 OK response with JSON data
func SuccessJSON(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(data)
}