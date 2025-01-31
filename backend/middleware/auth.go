package middleware

import (
	"context"
	"database/sql"

	"github.com/adamdgit/gotest/backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Checks User is logged in, for protected routes
func AuthLoggedIn(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the session, handles expired sessions automatically
		session, err := store.Get(c)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Check if user is logged in (e.g., session contains a user ID)
		userID := session.Get("id")
		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "You must be logged in to access this route",
			})
		}

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

		session_id := session.ID()

		// Get ID from cookie
		id, ok := session.Get("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Check database for user role
		stmt := "SELECT role FROM users WHERE id = ? AND session_id = ?"
		row := db.QueryRowContext(context.Background(), stmt, id, session_id)

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

// Reusable error handling function
// func HandleError(err error, message string, statusCode int, c *fiber.Ctx) error {
// 	if err != nil {
// 		return c.Status(statusCode).JSON(fiber.Map{
// 			"error": message,
// 		})
// 	}
// 	return nil
// }
