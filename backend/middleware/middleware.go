package middleware

import (
	"database/sql"
	"encoding/base64"
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
			log.Printf("AuthSessionIsValid ERR1: no access token")
			return c.Status(fiber.StatusUnauthorized).JSON(
				api.ErrAccessExpired,
			)
		}

		var access_expiration time.Time

		// incoming tokens must be decoded into binary (how they are stored in DB)
		decoded_access, err := base64.RawURLEncoding.DecodeString(access_token)
		if err != nil {
			log.Printf("AuthSessionIsValid ERR4: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}
		log.Printf("decoded token: %s", decoded_access)

		// check token exists in db
		err = db.QueryRow("SELECT access_expires FROM sessions WHERE access_token = ?", decoded_access).
			Scan(&access_expiration)
		if err != nil {
			log.Printf("AuthSessionIsValid ERR2: %s", err)
			return c.Status(fiber.StatusUnauthorized).JSON(
				api.ErrAccessExpired,
			)
		}

		// check token expiration
		if time.Now().UTC().After(access_expiration) {
			log.Printf("AuthSessionIsValid ERR3: expired token %s | %s", access_expiration, decoded_access)
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

// NOTE: AuthSessionIsValid() should always be run first,
// as this middleware has one function only, to verify roles
// valid session is assumed once reaching this middleware
func AuthUserHasRole(db *sql.DB, allowedRoles ...models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		access_token := c.Cookies("access_token")

		// incoming tokens must be decoded into binary (how they are stored in DB)
		decoded_access, err := base64.RawURLEncoding.DecodeString(access_token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Retrieve user role
		var user models.User

		err = db.QueryRow(`
            SELECT user_role
            FROM sessions
            WHERE access_token = ?
        `, decoded_access).Scan(&user.Role)
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
