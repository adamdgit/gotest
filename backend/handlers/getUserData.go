package handlers

import (
	"context"
	"database/sql"
	"log"

	"github.com/adamdgit/gotest/backend/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Deletes http-only cookie
func GetUserData(db *sql.DB, store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("context: %s", c.Request().Header.Header())
		// Retrieve the session if it exists
		session, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid session",
			})
		}

		// Get ID from cookie
		session_id := session.ID()
		user_id := session.Get("user_id")

		var user models.User

		// Check database for user data
		stmt := "SELECT email, role, firstname, lastname FROM users WHERE session_id = ? AND id = ?"

		row := db.QueryRowContext(context.Background(), stmt, session_id, user_id)
		err = row.Scan(&user.Email, &user.Role, &user.Firstname, &user.Lastname)
		if err == sql.ErrNoRows {
			log.Printf("error4")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Couldn't retrieve user data",
			})
		}

		// Success, return data as json
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"user": fiber.Map{"email": user.Email, "role": user.Role},
		})
	}
}
