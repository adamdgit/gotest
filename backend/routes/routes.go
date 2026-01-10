package routes

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/handlers"
	"github.com/adamdgit/gotest/backend/handlers/categories"
	"github.com/adamdgit/gotest/backend/handlers/products"
	"github.com/adamdgit/gotest/backend/handlers/users"
	"github.com/adamdgit/gotest/backend/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/oschwald/geoip2-golang"
)

func RegisterAPIRoutes(app *fiber.App, db *sql.DB, geoDb *geoip2.Reader) {
	// ----------------- PRODUCTS ----------------- //

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

	// ----------------- CATEGORIES ----------------- //

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

	// ----------------- OTHER ----------------- //

	// Refresh session using refresh token
	app.Get("/api/refresh", handlers.RefreshAccessToken(db))

	// Login, Logout, Register
	app.Post("/api/auth/login", handlers.Login(db, geoDb))
	app.Get("/api/auth/logout", handlers.Logout(db))
	app.Post("/api/auth/register", handlers.Register(db))

	// ----------------- USER ----------------- //

	// Gets users data based on session cookie
	app.Get("/api/user",
		middleware.AuthSessionIsValid(db),
		users.GetUserData(db),
	)

	// ----------------- ADMIN ----------------- //

	// Get user by ID
	app.Get("/api/admin/user/:id",
		middleware.AuthSessionIsValid(db),
		middleware.AuthUserHasRole(db, "admin"),
		users.GetUserDataByID(db),
	)

	// Get all user data by search query
	app.Get("/api/admin/users",
		middleware.AuthSessionIsValid(db),
		middleware.AuthUserHasRole(db, "admin"),
		users.GetUserByQuery(db),
	)
}
