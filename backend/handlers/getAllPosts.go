package handlers

import (
	"database/sql"
	"log"

	"github.com/adamdgit/gotest/backend/models"
	"github.com/gofiber/fiber/v2"
)

// get all posts
func GetAllPosts(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		stmt := "SELECT * FROM posts LIMIT 10"

		rows, err := db.Query(stmt)
		if err != nil {
			log.Printf("Error: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error retrieving from database",
			})
		}
		defer rows.Close()

		var posts []models.Post

		for rows.Next() {
			var post models.Post

			err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created_At, &post.Updated_At)
			if err != nil {
				log.Printf("Error: %s", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Error retrieving from database",
				})
			}
			posts = append(posts, post)
		}

		return c.JSON(posts)
	}
}
