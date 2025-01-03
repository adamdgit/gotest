package handlers

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Deletes http-only cookie
func Logout(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create session cookie
		session, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error creating session",
			})
		}

		// Destroy session
		if err := session.Destroy(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error logging out",
			})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
