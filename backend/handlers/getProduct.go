package handlers

import (
	"database/sql"
	"log"

	"github.com/adamdgit/gotest/backend/models"
	"github.com/gofiber/fiber/v2"
)

// Get post by provided id
func GetProduct(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		stmt := "SELECT * FROM products WHERE id = ?"

		row, err := db.Query(stmt, id)
		if err != nil {
			log.Printf("Error: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error retrieving from database",
			})
		}
		defer row.Close()

		var product models.Product

		err = row.Scan(&product.ID, &product.Name, &product.Brand, &product.Description, &product.Price, &product.Created_At, &product.Updated_At)
		if err != nil {
			log.Printf("Error: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error retrieving from database",
			})
		}

		return c.JSON(product)
	}
}
