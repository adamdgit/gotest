package routes

import (
	"database/sql"
	"log"

	"github.com/adamdgit/gotest/backend/queries"
	"github.com/gofiber/fiber/v2"
)

func RegisterAPIRoutes(app *fiber.App, db *sql.DB) {

	// Routes
	app.Get("/api/v1/posts", func(c *fiber.Ctx) error {
		posts, err := queries.GetAllPosts(app, db)
		if err != nil {
			log.Printf("error: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err,
			})
		}

		return c.JSON(posts)
	})

	app.Get("api/v1/posts/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return c.JSON("hello", id)
	})

}
