package routes

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app *fiber.App, db *sql.DB) {
	// Refresh session using refresh token
	app.Get("/api/auth/refresh", handlers.RefreshAccessToken(db))

	app.Post("/api/auth/login", handlers.Login(db))

	app.Get("/api/auth/logout", handlers.Logout(db))

	app.Post("/api/auth/register", handlers.Register(db))
}
