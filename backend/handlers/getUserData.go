package handlers

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

// Deletes http-only cookie
func GetUserData(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("context: %s", c.Request().Header.Header())
		// Retrieve the session if it exists
		sessionID := c.Cookies("access_token")

		if sessionID == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Couldn't retrieve user data",
			})
		}

		// TODO; select user_id from sessions table
		// Select user data from user table
		// Return data for user
		return nil
		// Success, return data as json
		// return c.Status(fiber.StatusOK).JSON(fiber.Map{
		// 	"user": fiber.Map{
		// 		"email":       user.Email,
		// 		"role":        user.Role,
		// 		"profile_url": user.Profile_URL,
		// 	},
		// })
	}
}
