package routes

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/handlers/products"
	"github.com/adamdgit/gotest/backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterProductRoutes(app *fiber.App, db *sql.DB) {
	app.Get("/api/v1/products",
		middleware.AuthSessionIsValid(db),
		products.GetProductList(db),
	)

	// Get product info by ID
	app.Get("/api/v1/product/:id",
		middleware.AuthSessionIsValid(db),
		products.GetProduct(db),
	)

	// Update product by ID
	app.Put("/api/v1/product/:id",
		middleware.AuthSessionIsValid(db),
		products.UpdateProduct(db),
	)

	// Get products for selected category
	app.Get("/api/v1/products/:categoryid",
		middleware.AuthSessionIsValid(db),
		products.GetProductsByCategory(db),
	)
}
