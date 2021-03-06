package handlers

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

// serverError returns all 5xx error messages
func serverError(c *fiber.Ctx, status int, error interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"error": error,
	})
}

// clientError returns all 5xx error messages
func clientError(c *fiber.Ctx, status int, error interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"error": error,
	})
}

// validationError returns all validation errors as  422
func validationError(c *fiber.Ctx, err error) error {
	errorsBody := make(map[string]string)

	for _, s := range strings.Split(err.Error(), ";") {
		values := strings.Split(s, ":")
		errorsBody[strings.Trim(values[0], " ")] = strings.Trim(values[1], " .")
	}

	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"errors": errorsBody,
	})
}

// validationDuplicateError returns all db duplicate errors as 422
func validationDuplicateError(c *fiber.Ctx, errors fiber.Map) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"errors": errors,
	})
}

// successResponse returns all 2xx responses
func successResponse(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"data": data,
	})
}
