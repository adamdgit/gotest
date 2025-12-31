package categories

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// JSON format from login body request
type CategoryReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Get post by provided id
func AddNewCategory(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CategoryReq

		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request params",
			})
		}

		name := req.Name
		description := req.Description

		var exists bool
		err = db.QueryRowContext(context.Background(),
			"SELECT EXISTS (SELECT 1 FROM categories WHERE name = ?)", name).Scan(&exists)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}
		if exists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Category already exists",
			})
		}

		_, err = db.Exec("INSERT INTO categories (name, description) VALUES (?, ?)",
			name, description)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "New category created successfully",
		})
	}
}
