package handlers

import (
	"database/sql"
	"log"

	"github.com/adamdgit/gotest/backend/models"
	"github.com/gofiber/fiber/v2"
)

// Get post by provided id
func GetCategories(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, name, description FROM categories")
		if err != nil {
			log.Printf("Error: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error connecting to database",
			})
		}
		defer rows.Close()

		var categories []models.Categories

		for rows.Next() {
			var category models.Categories
			err := rows.Scan(&category.ID, &category.Name, &category.Description)
			if err != nil {
				log.Printf("Error: %s", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Error retrieving from database",
				})
			}
			categories = append(categories, category)
		}

		return c.JSON(categories)
	}
}
