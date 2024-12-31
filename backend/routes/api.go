package routes

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/handlers"
	"github.com/adamdgit/gotest/backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterAPIRoutes(app *fiber.App, db *sql.DB) {
	// Load JWT secret from .env
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	// secret := os.Getenv("JWT_SECRET")
	authRequired := middleware.NewAuthMiddleware("secret")
	// NOTE: Add authRequired to routes that need auth protection

	// Post Routes
	app.Get("/api/v1/posts", handlers.GetAllPosts(db))
	app.Get("/api/v1/posts/:id", handlers.GetPostById(db))

	// Login & Register routes
	app.Post("/api/auth/login", handlers.Login(db))
	app.Post("/api/auth/logout", authRequired, handlers.Logout) // clears session cookie, no db required
	app.Post("/api/auth/register", handlers.Register(db))

}
