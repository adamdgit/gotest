package handlers

import (
	"database/sql"
	"log"

	"github.com/adamdgit/gotest/backend/models"
	"github.com/gofiber/fiber/v2"
)

// Get post by provided id
func GetPostById(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		stmt := "SELECT * FROM posts WHERE post.id = ?"

		row, err := db.Query(stmt, id)
		if err != nil {
			log.Printf("Error: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error retrieving from database",
			})
		}
		defer row.Close()

		var post models.Post

		err = row.Scan(&post.ID, &post.Title, &post.Content, &post.Created_At, &post.Updated_At)
		if err != nil {
			log.Printf("Error: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error retrieving from database",
			})
		}

		return c.JSON(post)
	}
}
