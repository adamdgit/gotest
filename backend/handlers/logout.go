package handlers

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

// Deletes http-only cookie
func Logout(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		refreshToken := c.Cookies("refresh_token")

		if refreshToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No session found",
			})
		}

		// Delete session from db
		_, err := db.Exec("DELETE FROM sessions WHERE refresh_token = ?", refreshToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Error logging out",
			})
		}

		// Delete cookies on logout
		c.Cookie(&fiber.Cookie{
			Name:    "access_token",
			Value:   "",
			Expires: time.Now().Add(-time.Hour),
		})

		c.Cookie(&fiber.Cookie{
			Name:    "refresh_token",
			Value:   "",
			Expires: time.Now().Add(-time.Hour),
		})

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"meessage": "Logged out successfully",
		})
	}
}
