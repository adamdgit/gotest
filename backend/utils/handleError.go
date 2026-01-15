package utils

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// HandleAPIError logs optionally and sends a JSON error response
func HandleAPIError(
	c *fiber.Ctx,
	err error,
	status int,
	response any,
	logMsg *string,
) error {
	if err == nil {
		return nil
	}

	if logMsg != nil {
		log.Printf("%s: %v", *logMsg, err)
	}

	return c.Status(status).JSON(response)
}
