package routes

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/handlers"
	"github.com/adamdgit/gotest/backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterAPIRoutes(app *fiber.App, db *sql.DB) {
	// NOTE: Protect routes by adding AuthLoggedIn or AuthIsAdmin
	// AuthLoggedIn, users must be logged in to access API
	// AuthIsAdmin, users must have admin role to access API

	app.Get("/api/v1/products", middleware.AuthSessionIsValid(db), handlers.GetProductList(db))
	app.Get("/api/v1/product/:id",
		middleware.AuthSessionIsValid(db),
		handlers.GetProduct(db),
	)
	app.Put("/api/v1/product/:id",
		middleware.AuthIsAdmin(db),
		handlers.UpdateProduct(db),
	)

	// Login, Logout, Register
	app.Post("/api/auth/login", handlers.Login(db))
	app.Get("/api/auth/logout", handlers.Logout(db))
	app.Post("/api/auth/register", handlers.Register(db))

	// Validate users session after login, before redirect to main app
	app.Get("/api/auth/getUser", handlers.GetUserData(db))
}
