package routes

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func RegisterAPIRoutes(app *fiber.App, db *sql.DB, store *session.Store) {
	// NOTE: Protect routes by adding AuthLoggedIn or AuthIsAdmin
	// AuthLoggedIn, users must be logged in to access API
	// AuthIsAdmin, users must have admin role to access API
	app.Get("/api/v1/products", handlers.GetProductList(db))
	app.Get("/api/v1/product/:id", handlers.GetProduct(db))

	// Login, Logout, Register
	app.Post("/api/auth/login", handlers.Login(db, store))
	app.Get("/api/auth/logout", handlers.Logout(db, store))
	app.Post("/api/auth/register", handlers.Register(db))

	// Validate users session after login, before redirect to main app
	app.Get("/api/auth/getUser", handlers.GetUserData(db, store))
}
