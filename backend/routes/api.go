package routes

import (
	"github.com/adamdgit/gotest/backend/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterAPIRoutes(app *fiber.App) {

	// Post Routes
	app.Get("/api/v1/posts", handlers.GetAllPosts)
	app.Get("/api/v1/posts/:id", handlers.GetPostById)

	// Login & Register routes
	app.Post("/api/auth/login", handlers.Login)
	app.Post("/api/auth/register", handlers.Register)

}
