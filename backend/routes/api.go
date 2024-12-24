package routes

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/models"
	"github.com/gofiber/fiber/v2"
)

func RegisterAPIRoutes(app *fiber.App, db *sql.DB) {
	// Routes
	app.Get("/api/v1/posts", func(c *fiber.Ctx) error {
		stmt := "SELECT * FROM posts"

		rows, err := db.Query(stmt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch posts",
			})
		}
		defer rows.Close()

		var posts []models.Post
		for rows.Next() {
			var post models.Post

			err := rows.Scan(&post.ID, &post.Title, &post.Content)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to scan post",
				})
			}
			posts = append(posts, post)
		}

		return c.JSON(posts)
	})
}
