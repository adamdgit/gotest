package handlers

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Deletes http-only cookie
func Logout(db *sql.DB, store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve session cookie if it exists
		session, err := store.Get(c)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		user_id := session.Get("user_id")
		// handles null types for db
		var nullSessionID sql.NullString
		nullSessionID.Valid = false

		// Remove session from user in db
		_, err = db.Exec("UPDATE users SET session_id = ? WHERE id = ?", nullSessionID, user_id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error logging out",
			})
		}

		// Destroy session
		if err := session.Destroy(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error logging out",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Logged out successfully",
		})
	}
}
