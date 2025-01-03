package routes

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/handlers"
	"github.com/adamdgit/gotest/backend/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func RegisterAPIRoutes(app *fiber.App, db *sql.DB, store *session.Store) {
	// NOTE: Protect routes by adding auth middleware:
	// middleware.AuthMiddleware(store)
	app.Get("/api/v1/posts", middleware.AuthMiddleware(store), handlers.GetAllPosts(db))
	app.Get("/api/v1/posts/:id", handlers.GetPostById(db))

	// Login, Logout, Register
	app.Post("/api/auth/login", handlers.Login(db, store))
	app.Post("/api/auth/logout", handlers.Logout(store))
	app.Post("/api/auth/register", handlers.Register(db))
}
