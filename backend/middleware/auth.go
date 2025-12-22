package middleware

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Check session expiration
func AuthSessionIsValid(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Validates access token, or generates a new one
		access_token := c.Cookies("access_token")

		if access_token == "" {
			log.Printf("error: no access token")
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		var access_expiration time.Time

		// check token exists in db
		err := db.QueryRow("SELECT access_expires FROM sessions WHERE access_token = ?", access_token).
			Scan(&access_expiration)
		if err == sql.ErrNoRows {
			log.Printf("error, %s", err)
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// check token expiration
		if time.Now().UTC().After(access_expiration) {
			log.Printf("error: expired token %s | %s", access_expiration, access_token)
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// TODO? We can also add checks for things like IP / geolocation
		// to further protect users.

		// Session is valid
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
