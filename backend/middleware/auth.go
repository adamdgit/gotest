package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Checks users session cookie for protecting routes
func AuthMiddleware(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the session
		session, err := store.Get(c)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Check if user is logged in (e.g., session contains a user ID)
		userID := session.Get("id")
		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized for this route",
			})
		}

		return c.Next()
	}
}
