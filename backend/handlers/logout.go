package handlers

import (
	"database/sql"
	"encoding/base64"
	"time"

	"github.com/adamdgit/gotest/backend/api"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

// Deletes http-only cookie
func Logout(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		refresh_token := c.Cookies("refresh_token")

		if refresh_token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				api.ErrInvalidRequest,
			)
		}

		// incoming tokens must be decoded into binary (how they are stored in DB)
		decoded_refresh, err := base64.RawURLEncoding.DecodeString(refresh_token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Delete session from db
		_, err = db.Exec("DELETE FROM sessions WHERE refresh_token = ?", decoded_refresh)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				api.ErrInternalServer,
			)
		}

		// Delete cookies on logout
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    "",
			HTTPOnly: true,
			Secure:   false,
			SameSite: "None",
			Expires:  time.Unix(0, 0),
		})

		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    "",
			HTTPOnly: true,
			Secure:   false,
			SameSite: "None",
			Expires:  time.Unix(0, 0),
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Logged out successfully",
		})
	}
}
