package routes

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/handlers"
	"github.com/adamdgit/gotest/backend/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/oschwald/geoip2-golang"
)

func RegisterAPIRoutes(app *fiber.App, db *sql.DB, geoDb *geoip2.Reader) {
	// Get all Products
	app.Get("/api/v1/products",
		middleware.AuthSessionIsValid(db),
		handlers.GetProductList(db),
	)
	// Get product info by ID
	app.Get("/api/v1/product/:id",
		middleware.AuthSessionIsValid(db),
		handlers.GetProduct(db),
	)
	// Update product by ID
	app.Put("/api/v1/product/:id",
		middleware.AuthIsAdmin(db),
		handlers.UpdateProduct(db),
	)

	// Get all categories
	app.Get("/api/v1/categories",
		middleware.AuthSessionIsValid(db),
		handlers.GetCategories(db),
	)
	// Add new categories
	app.Put("/api/v1/categories",
		middleware.AuthSessionIsValid(db),
		handlers.AddNewCategory(db),
	)
	// Get products for selected category
	app.Get("/api/v1/products/:categoryid",
		middleware.AuthSessionIsValid(db),
		handlers.GetProductsByCategory(db),
	)

	// Login, Logout, Register
	app.Post("/api/auth/login", handlers.Login(db, geoDb))
	app.Get("/api/auth/logout", handlers.Logout(db))
	app.Post("/api/auth/register", handlers.Register(db))

	// Gets current users data
	app.Get("/api/user",
		middleware.AuthSessionIsValid(db),
		handlers.GetUserData(db),
	)

	app.Get("/api/refresh", handlers.RefreshAccessToken(db))
}
