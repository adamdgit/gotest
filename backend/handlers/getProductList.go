package handlers

import (
	"database/sql"
	"log"

	"github.com/adamdgit/gotest/backend/models"
	"github.com/gofiber/fiber/v2"
)

// get all posts
func GetProductList(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		stmt := "SELECT * FROM products LIMIT 20"

		rows, err := db.Query(stmt)
		if err != nil {
			log.Printf("Error: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error retrieving from database",
			})
		}
		defer rows.Close()

		var products []models.Product

		for rows.Next() {
			var product models.Product

			err := rows.Scan(&product.ID, &product.Name, &product.Brand, &product.Description, &product.Price, &product.Created_At, &product.Updated_At)
			if err != nil {
				log.Printf("Error: %s", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Error retrieving from database",
				})
			}
			products = append(products, product)
		}

		// return c.Render("test", fiber.Map{
		// 	"products": products,
		// })
		return c.JSON(products)
	}
}
