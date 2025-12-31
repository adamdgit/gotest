package middleware

import (
	"database/sql"
	"log"
	"slices"
	"time"

	"github.com/adamdgit/gotest/backend/api"
	"github.com/adamdgit/gotest/backend/models"
	"github.com/gofiber/fiber/v2"
)

// Check session expiration
func AuthSessionIsValid(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Validates access token, or generates a new one
		access_token := c.Cookies("access_token")

		if access_token == "" {
			log.Printf("AuthSessionIsValid ERR: no access token")
			return c.Status(fiber.StatusUnauthorized).JSON(
				api.ErrAccessExpired,
			)
		}

		var access_expiration time.Time

		// check token exists in db
		err := db.QueryRow("SELECT access_expires FROM sessions WHERE access_token = ?", access_token).
			Scan(&access_expiration)
		if err != nil {
			log.Printf("AuthSessionIsValid ERR, %s", err)
			return c.Status(fiber.StatusUnauthorized).JSON(
				api.ErrAccessExpired,
			)
		}

		// check token expiration
		if time.Now().UTC().After(access_expiration) {
			log.Printf("AuthSessionIsValid ERR: expired token %s | %s", access_expiration, access_token)
			return c.Status(fiber.StatusUnauthorized).JSON(
				api.ErrAccessExpired,
			)
		}

		// TODO? We can also add checks for things like IP / geolocation
		// to further protect users.

		// Session is valid
		return c.Next()
	}
}

// AuthIsAdmin middleware ensures the user has admin privileges
// NOTE: AuthSessionIsValid() should always be run first,
// so we don't need to check the session is valid again
func AuthUserHasRole(db *sql.DB, allowedRoles ...models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessToken := c.Cookies("access_token")

		// Retrieve user role
		var user models.User

		err := db.QueryRow(`
            SELECT u.id, u.role
            FROM users u
            INNER JOIN sessions s ON u.id = s.user_id
            WHERE s.access_token = ?
        `, accessToken).Scan(&user.ID, &user.Role)
		if err != nil {
			log.Printf("error: get user role %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// If role isn't allowed, throw error
		if !slices.Contains(allowedRoles, user.Role) {
			return c.Status(fiber.StatusForbidden).JSON(
				api.ErrForbidden,
			)
		}

		// Role is valid, continue
		return c.Next()
	}
}
