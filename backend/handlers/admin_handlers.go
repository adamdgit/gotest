package handlers

import (
	"database/sql"
	"log"
	"strings"

	"github.com/adamdgit/gotest/backend/api"
	"github.com/adamdgit/gotest/backend/services"
	"github.com/gofiber/fiber/v2"
)

// Get all user data by email, ADMIN access only
func GetUserByQuery(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := strings.TrimSpace(c.Query("search"))
		if query == "" {
			return c.Status(fiber.StatusBadRequest).JSON(
				api.ErrInvalidBody,
			)
		}

		// Handle any other query validation here as necessary

		users, err := services.DB_GetUserDataByQuery(db, query)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		return c.Status(fiber.StatusOK).JSON(users)
	}
}

func GetUserDataByID(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("userid")

		user, err := services.DB_GetUserDataByID(db, id)
		if err != nil {
			log.Printf("Error: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		return c.Status(fiber.StatusOK).JSON(user)
	}
}
