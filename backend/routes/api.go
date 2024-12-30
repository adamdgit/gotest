package routes

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterAPIRoutes(app *fiber.App, db *sql.DB) {
	// Load JWT secret from .env
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	// secret := os.Getenv("JWT_SECRET")
	// jwt := middleware.NewAuthMiddleware(secret)
	// NOTE: Add jwt to routes that need auth protection

	// Post Routes
	app.Get("/api/v1/posts", handlers.GetAllPosts(db))
	app.Get("/api/v1/posts/:id", handlers.GetPostById(db))

	// Login & Register routes
	app.Post("/api/auth/login", handlers.Login(db))
	app.Post("/api/auth/register", handlers.Register(db))

}
