package handlers

import (
	"log"

	"github.com/adamdgit/gotest/backend/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Deletes http-only cookie
func GetUser(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Print("context: ", c)
		// Retrieve the session if it exists
		session, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid session",
			})
		}

		// User session token must contain a role or email, otherwise something is wrong
		role, ok := session.Get("role").(models.UserRole)
		log.Print("role: ", string(role))
		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		email, ok := session.Get("email").(string)
		log.Print("role: ", models.UserRole(email))
		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Success, return data as json
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"user": fiber.Map{"email": email, "role": role},
		})
	}
}
