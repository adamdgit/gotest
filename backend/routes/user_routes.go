package routes

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/handlers"
	"github.com/adamdgit/gotest/backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App, db *sql.DB) {
	// Gets users data based on session cookie
	app.Get("/api/user",
		middleware.AuthSessionIsValid(db),
		handlers.GetUserData(db),
	)
}
