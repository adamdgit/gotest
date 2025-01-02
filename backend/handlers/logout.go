package handlers

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

// Deletes http-only cookie
func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(-time.Hour),
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(-time.Hour),
	})

	return c.SendStatus(fiber.StatusOK)
}
