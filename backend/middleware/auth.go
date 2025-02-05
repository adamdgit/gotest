package middleware

import (
	"context"
	"database/sql"

	"github.com/adamdgit/gotest/backend/models"
	"github.com/adamdgit/gotest/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
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
func AuthIsAdmin(db *sql.DB, store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the session if it exists
		session, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Get ID from cookie
		id, ok := session.Get("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Check database for user role
		stmt := "SELECT role FROM users WHERE id = ?"
		row := db.QueryRowContext(context.Background(), stmt, id)

		var role models.UserRole

		// If ErrNoRows user has provided invalid login details
		// else we need to check password is valid
		err = row.Scan(&role)
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		if role != models.Admin {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		return c.Next()
	}
}
