package routes

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/handlers/categories"
	"github.com/adamdgit/gotest/backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterCategoryRoutes(app *fiber.App, db *sql.DB) {
	// Get all categories
	app.Get("/api/v1/categories",
		middleware.AuthSessionIsValid(db),
		categories.GetCategories(db),
	)

	// Add new categories
	app.Put("/api/v1/categories",
		middleware.AuthSessionIsValid(db),
		categories.AddNewCategory(db),
	)
}
