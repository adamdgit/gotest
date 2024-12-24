package routes

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/adamdgit/gotest/backend/queries"
	"github.com/gofiber/fiber/v2"
)

func RegisterAPIRoutes(app *fiber.App, db *sql.DB) {
	apiVersion := os.Getenv("API_VERSION")

	// Routes
	app.Get(fmt.Sprintf("/api/%s/posts", apiVersion), func(c *fiber.Ctx) error {
		posts, err := queries.GetAllPosts(app, db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err,
			})
		}

		return c.JSON(posts)
	})

	app.Get(fmt.Sprintf("/api/%s/posts/:id", apiVersion), func(c *fiber.Ctx) error {
		id := c.Params("id")
		return c.JSON("hello", id)
	})

}
