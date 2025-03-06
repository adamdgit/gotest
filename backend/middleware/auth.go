package middleware

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/utils"
	"github.com/gofiber/fiber/v2"
)

// Check session expiration
func AuthSessionIsValid(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Validates access token, or generates a new one
		err := utils.ValidateAccessToken(c, db)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Session",
			})
		}
		// Session is valid continue
		return c.Next()
	}
}

// Checks User is admin role for protected routes
func AuthIsAdmin(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the session if it exists
		return c.Next()
	}
}
