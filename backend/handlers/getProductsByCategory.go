package handlers

import (
	"database/sql"
	"log"

	"github.com/adamdgit/gotest/backend/models"
	"github.com/gofiber/fiber/v2"
)

// Get post by provided id
func GetProductsByCategory(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("categoryid")

		rows, err := db.Query("SELECT * FROM products WHERE category = ? LIMIT 20", id)
		if err != nil {
			log.Printf("Error: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error connecting to database",
			})
		}
		defer rows.Close()

		var products []models.Product

		for rows.Next() {
			var product models.Product
			err := rows.Scan(&product.ID, &product.Name, &product.Brand, &product.Description, &product.Price, &product.Category, &product.Created_At, &product.Updated_At)
			if err != nil {
				log.Printf("Error: %s", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Error retrieving from database",
				})
			}
			products = append(products, product)
		}

		return c.JSON(products)
	}
}
