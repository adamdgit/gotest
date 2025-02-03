package utils

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// HandleError handles different types of errors and sends a json response
func HandleError(c *fiber.Ctx, err error, message string) error {
	if err == nil {
		return nil // No error, so just return nil
	}

	// Check for specific error types
	switch err {
	case nil:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": message,
		})

	case sql.ErrNoRows:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": message,
		})

	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	}
}
